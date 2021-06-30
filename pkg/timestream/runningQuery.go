package timestream

import (
	"context"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/timestream-datasource/pkg/models"
)

type openQuery struct {
	nextToken string
	query     models.QueryModel
	runner    queryRunner
}

func (ds *openQuery) doStream(ctx context.Context, sender *backend.StreamSender) error {
	timer := time.NewTimer(time.Second * 10)

	for {
		select {
		case <-ctx.Done():
			backend.Logger.Info("stop streaming (context canceled)")
			return nil
		case <-timer.C:
			backend.Logger.Info("TODO... run query!", "token", ds.nextToken)
		}
	}
}
