package timestream

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/timestream-datasource/pkg/models"
)

// QueryResultToDataFrame creates a DataFrame from query results
func QueryResultToDataFrame(res *timestreamquery.QueryOutput) (dr backend.DataResponse) {
	var timeseriesColumn *fieldBuilder
	dr = backend.DataResponse{}
	notices := []data.Notice{}
	builders := []fieldBuilder{}
	cellParsingError := false
	length := len(res.Rows)

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
			if timeseriesColumn != nil {
				dr.Error = fmt.Errorf("Multiple timeseries columns are not supported")
				return
			}
			timeseriesColumn = b
		}
		builders = append(builders, *b)
	}

	// Each row is a new series
	if timeseriesColumn != nil {
		for _, series := range res.Rows {
			tv := series.Data[timeseriesColumn.columnIdx].TimeSeriesValue
			if tv == nil {
				dr.Error = fmt.Errorf("Expecting timeseries colum at: %d", timeseriesColumn.columnIdx)
				return
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
				v, _ := timeseriesColumn.parser(tv[i].Value)
				tf.Set(i, t)
				vf.Set(i, v)
			}

			// Add the series as a frame
			dr.Frames = append(dr.Frames, data.NewFrame("", tf, vf))
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
				if row == nil {
					continue
				}
				v, err := builder.parser(row.Data[builder.columnIdx])
				if err != nil {
					if !cellParsingError {
						notices = append(notices, data.Notice{
							Severity: data.NoticeSeverityError,
							Text:     fmt.Sprintf("Error parsing: row:%d, colum:%d", i, builder.columnIdx),
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

		frame := data.NewFrame("", // No name
			fields...,
		)
		dr.Frames = append(dr.Frames, frame)
	}

	meta := &models.TimestreamCustomMeta{
		QueryID:   aws.StringValue(res.QueryId),
		NextToken: aws.StringValue(res.NextToken),
		HasSeries: timeseriesColumn != nil,
	}

	// At least one empty result
	if len(dr.Frames) < 1 {
		dr.Frames = append(dr.Frames, data.NewFrame(""))
	}

	// Attach all notices to the first response
	if len(notices) > 0 {
		dr.Frames[0].AppendNotices(notices...)
	}
	if dr.Frames[0].Meta == nil {
		dr.Frames[0].Meta = &data.FrameMeta{}
	}
	dr.Frames[0].Meta.Custom = meta
	return
}
