package timestream

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/errorsource"
	"github.com/grafana/timestream-datasource/pkg/models"
)

// QueryResultToDataFrame creates a DataFrame from query results
func QueryResultToDataFrame(res *timestreamquery.QueryOutput, format models.FormatQueryOption) backend.DataResponse {
	dr := backend.DataResponse{}
	notices := []data.Notice{}
	builders := []*fieldBuilder{}
	timeseriesColumns := []*fieldBuilder{}
	cellParsingError := false
	length := len(res.Rows)
	hasTimeseries := false

	// Inspect the column structure
	for index, columnMeta := range res.ColumnInfo {
		b, err := getFieldBuilder(columnMeta.Type)
		if err != nil {
			notices = append(notices, data.Notice{
				Severity: data.NoticeSeverityWarning,
				Text:     err.Error(),
			})
			continue
		}
		b.columnIdx = index
		b.name = *columnMeta.Name
		if b.timeseries {
			timeseriesColumns = append(timeseriesColumns, b)
			hasTimeseries = true
		} else {
			builders = append(builders, b)
		}
	}

	if hasTimeseries {
		// Each row is a new series
		for _, timeseriesColumn := range timeseriesColumns {
			for _, series := range res.Rows {
				tv := series.Data[timeseriesColumn.columnIdx].TimeSeriesValue
				nv := series.Data[timeseriesColumn.columnIdx].NullValue
				isNullDataPoint := nv != nil && *nv
				if tv == nil && !isNullDataPoint {
					return errorsource.Response(errorsource.PluginError(
						fmt.Errorf("expecting timeseries column at: %d", timeseriesColumn.columnIdx), false))
				}

				length := len(tv)
				tf := data.NewFieldFromFieldType(data.FieldTypeTime, length)
				vf := data.NewFieldFromFieldType(timeseriesColumn.fieldType, length)
				tf.Name = "time"
				vf.Name = timeseriesColumn.name
				vf.Labels = data.Labels{}
				for _, builder := range builders {
					val := series.Data[builder.columnIdx].ScalarValue
					if !builder.timeseries && val != nil {
						vf.Labels[builder.name] = *val
					}
				}

				for i := 0; i < length; i++ {
					t, _ := time.Parse("2006-01-02 15:04:05.99999999", *tv[i].Time)
					v, _ := timeseriesColumn.parser(*tv[i].Value)
					tf.Set(i, t)
					vf.Set(i, v)
				}

				// Add the series as a frame
				dr.Frames = append(dr.Frames, data.NewFrame("", tf, vf))
			}
		}
	} else {
		fields := []*data.Field{}
		for _, builder := range builders {
			field := data.NewFieldFromFieldType(builder.fieldType, length)
			field.Name = builder.name
			if builder.config != nil {
				field.Config = builder.config
			}
			for i := 0; i < length; i++ {
				row := res.Rows[i]
				v, err := builder.parser(row.Data[builder.columnIdx])
				if err != nil {
					if !cellParsingError {
						notices = append(notices, data.Notice{
							Severity: data.NoticeSeverityError,
							Text:     fmt.Sprintf("Error parsing: row:%d, column:%d", i, builder.columnIdx),
						})
					}
					cellParsingError = true
				} else if v != nil {
					// Convert json values to strings
					if builder.asJSON {
						bytes, err := json.Marshal(v)
						if err != nil {
							v = fmt.Sprintf("ERROR: %s", err.Error())
						} else {
							v = string(bytes)
						}
					}
					field.Set(i, v)
				}
			}
			fields = append(fields, field)
		}

		frame := data.NewFrame("", fields...)

		if length > 0 && format == models.FormatOptionTimeSeries {
			if frame.TimeSeriesSchema().Type == data.TimeSeriesTypeLong {
				var err error
				frame, err = data.LongToWide(frame, &data.FillMissing{
					Mode: data.FillModeNull,
				})
				if err != nil {
					return errorsource.Response(errorsource.PluginError(fmt.Errorf("error formatting as timeseries: %s", err), false))
				}
			}
		}
		dr.Frames = append(dr.Frames, frame)
	}

	meta := &models.TimestreamCustomMeta{
		HasSeries: hasTimeseries,
	}
	if res.QueryId != nil {
		meta.QueryID = *res.QueryId
	}
	if res.NextToken != nil {
		meta.NextToken = *res.NextToken
	}

	// At least one empty result
	if len(dr.Frames) == 0 {
		dr.Frames = data.Frames{data.NewFrame("")}
	}

	// Attach all notices to the first response
	if len(notices) > 0 {
		dr.Frames[0].AppendNotices(notices...)
	}
	if dr.Frames[0].Meta == nil {
		dr.Frames[0].Meta = &data.FrameMeta{}
	}
	dr.Frames[0].Meta.Custom = meta
	return dr
}
