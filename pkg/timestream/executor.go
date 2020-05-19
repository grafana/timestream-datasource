package timestream

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/timestream-datasource/pkg/models"
)

func ExecuteQuery(ctx context.Context, query models.QueryModel, runner queryRunner) (dr backend.DataResponse) {
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
	if err != nil {
		dr.Error = err
		return
	}

	frame, err := QueryResultToDataFrame(output)
	dr.Frames = append(dr.Frames, frame)
	dr.Error = err
	return
}
