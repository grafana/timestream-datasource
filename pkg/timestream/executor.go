package timestream

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	gaws "github.com/grafana/timestream-datasource/pkg/common/aws"
	"github.com/grafana/timestream-datasource/pkg/models"
)

// ExecuteQuery -- run a query
func ExecuteQuery(ctx context.Context, query models.QueryModel, runner queryRunner, settings gaws.DatasourceSettings) (dr backend.DataResponse) {
	dr = backend.DataResponse{}

	raw, err := Interpolate(query, settings)
	if err != nil {
		dr.Error = err
		return
	}
	input := &timestreamquery.QueryInput{
		QueryString: aws.String(raw),
	}

	if len(query.NextToken) > 0 {
		input.NextToken = aws.String(query.NextToken)
		backend.Logger.Info("running continue query", "token", query.NextToken)
	}

	start := time.Now().UnixNano() / 1000000
	output, err := runner.runQuery(ctx, input)
	finish := time.Now().UnixNano() / 1000000
	if err == nil {
		dr = QueryResultToDataFrame(output)
	} else {
		dr.Error = err
	}

	// Needs a frame for the metadata... even if just error
	if len(dr.Frames) < 1 {
		dr.Frames = append(dr.Frames, data.NewFrame(""))
	}
	frame := dr.Frames[0]
	if frame.Meta == nil {
		frame.SetMeta(&data.FrameMeta{})
	}
	frame.Meta.ExecutedQueryString = raw

	if frame.Meta.Custom == nil {
		frame.Meta.Custom = &models.TimestreamCustomMeta{}
	}

	// Apply the timing info
	meta := frame.Meta.Custom.(*models.TimestreamCustomMeta)
	if meta.NextToken == "" {
		meta.FinishTime = finish
	}
	if input.NextToken == nil {
		meta.StartTime = start
	}
	return
}
