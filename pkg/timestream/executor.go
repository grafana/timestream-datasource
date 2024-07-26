package timestream

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/errorsource"
	"github.com/grafana/timestream-datasource/pkg/models"
)

// ExecuteQuery -- run a query
func ExecuteQuery(ctx context.Context, query models.QueryModel, runner queryRunner, settings models.DatasourceSettings) backend.DataResponse {

	raw, err := Interpolate(query, settings)
	if err != nil {
		return errorsource.Response(err)
	}
	input := &timestreamquery.QueryInput{
		QueryString: aws.String(raw),
	}

	if query.NextToken != "" {
		input.NextToken = aws.String(query.NextToken)
		backend.Logger.Info("running continue query", "token", query.NextToken)
	}

	start := time.Now().UnixMilli()
	output, err := runner.runQuery(ctx, input)
	if query.WaitForResult && output.NextToken != nil && err == nil {
		for output.NextToken != nil {
			newPageInput := *input
			newPageInput.NextToken = output.NextToken
			newPageOutput, newPageErr := runner.runQuery(ctx, &newPageInput)
			if newPageErr != nil {
				err = newPageErr
				output.NextToken = nil
				continue
			}
			output.Rows = append(output.Rows, newPageOutput.Rows...)
			output.NextToken = newPageOutput.NextToken
		}
	}

	dr := backend.DataResponse{}
	if err == nil {
		dr = QueryResultToDataFrame(output, query.Format)
	} else {
		// override: false here because runQuery may return a PluginError
		dr = errorsource.Response(errorsource.DownstreamError(err, false))
	}
	finish := time.Now().UnixMilli()

	// Needs a frame for the metadata... even if just error
	if len(dr.Frames) == 0 {
		dr.Frames = data.Frames{data.NewFrame("")}
	}
	frame := dr.Frames[0]
	if frame.Meta == nil {
		frame.SetMeta(&data.FrameMeta{})
	}
	frame.Meta.ExecutedQueryString = raw

	if frame.Meta.Custom == nil {
		frame.Meta.Custom = &models.TimestreamCustomMeta{}
	}
	if output != nil && output.QueryStatus != nil {
		c := frame.Meta.Custom.(*models.TimestreamCustomMeta)
		c.Status = output.QueryStatus
	}

	// Apply the timing info
	meta := frame.Meta.Custom.(*models.TimestreamCustomMeta)
	if meta.NextToken == "" {
		meta.FinishTime = finish
	}
	if input.NextToken == nil {
		meta.StartTime = start
	}
	return dr
}
