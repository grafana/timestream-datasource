package timestream

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/timestream-datasource/pkg/models"
)

// QueryResultMetaStat from https://github.com/grafana/grafana/blob/277aee864269753d45a2a2e725998aac59f592e9/packages/grafana-data/src/types/data.ts#L48
type QueryResultMetaStat struct {
	DisplayName string `json:"displayName"`

	Value float64 `json:"value"`

	Unit string `json:"unit"`
}

// ExecuteQuery -- run a query
func ExecuteQuery(ctx context.Context, query models.QueryModel, runner queryRunner) (dr backend.DataResponse) {
	start := time.Now()
	dr = backend.DataResponse{}

	raw, err := Interpolate(query)
	if err != nil {
		dr.Error = err
		return
	}

	backend.Logger.Info("running query", "query", raw)

	input := &timestreamquery.QueryInput{
		QueryString: aws.String(raw),
	}

	output, err := runner.runQuery(ctx, input)
	if err == nil {
		dr = QueryResultToDataFrame(output)
	} else {
		dr.Error = err
	}

	if len(dr.Frames) < 1 {
		dr.Frames = append(dr.Frames, data.NewFrame(""))
	}
	frame := dr.Frames[0]

	// Add all the stats
	meta := make(map[string]interface{})
	if output != nil {
		meta["queryId"] = output.QueryId
		if output.NextToken != nil {
			meta["nextToken"] = output.NextToken
		}
	}

	rows, _ := frame.RowLen()
	if output.NextToken != nil {
		if rows > 0 {
			frame.AppendNotices(data.Notice{
				Severity: data.NoticeSeverityWarning,
				Text:     "More results not yet supported",
			})
		} else if err == nil {
			err = fmt.Errorf("Slow query not yet supported")
		}
	}

	if frame.Meta == nil {
		frame.Meta = &data.FrameMeta{}
	}

	frame.Meta.Custom = meta
	frame.Meta.ExecutedQueryString = raw

	stats := make([]QueryResultMetaStat, 1)
	stats[0] = QueryResultMetaStat{
		DisplayName: "Execution time",
		Value:       float64(time.Since(start).Milliseconds()),
		Unit:        "ms",
	}
	frame.Meta.Stats = stats
	return
}
