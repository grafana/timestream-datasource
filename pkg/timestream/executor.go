package timestream

import (
	"context"
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

	// Add all the stats
	meta := make(map[string]interface{})
	meta["executedQuery"] = raw

	var frame *data.Frame
	output, err := runner.runQuery(ctx, input)
	if err == nil {
		frame, err = QueryResultToDataFrame(output)
		// if err == nil {
		// 	// make columns into tags
		// 	frame, err = data.LongToWide(frame, &data.FillMissing{
		// 		Mode: data.FillModeNull,
		// 	})
		// }

		meta["queryId"] = output.QueryId
		meta["nextToken"] = output.NextToken
	}
	if frame == nil {
		frame = data.NewFrame("")
	}

	stats := make([]QueryResultMetaStat, 1)
	stats[0] = QueryResultMetaStat{
		DisplayName: "Execution time",
		Value:       float64(time.Since(start).Milliseconds()),
		Unit:        "ms",
	}

	frame.Meta = &data.FrameMeta{
		Custom: meta,
		Stats:  stats,
	}

	dr.Frames = append(dr.Frames, frame)
	dr.Error = err
	return
}
