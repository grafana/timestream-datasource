package timestream

import (
	"context"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/timestream-datasource/pkg/models"
)

type openQuery struct {
	queryId string
	query   *models.QueryModel
	ds      *timestreamDS
}

func (q *openQuery) doStream(ctx context.Context, sender *backend.StreamSender) error {
	backend.Logger.Info("Starting stream for", "queryId", q.queryId, "token", q.query.NextToken)
	timer := time.NewTimer(time.Second * 2)

	for {
		select {
		case <-ctx.Done():
			backend.Logger.Info("stop streaming (context canceled)")
			return nil
		case <-timer.C:
			backend.Logger.Info("TODO... run query!", "token", q.query.NextToken)

			dr, nextToken := ExecuteQuery(ctx, *q.query, q.ds.Runner, q.ds.Settings)
			if dr.Error != nil {
				backend.Logger.Error("error running streaming query", "query", q.queryId, "err", dr.Error)
				return nil
			}

			// Send each frame
			for _, frame := range dr.Frames {
				err := sender.SendFrame(frame, data.IncludeAll)
				if err != nil {
					log.DefaultLogger.Error(fmt.Sprintf("unable to send message: %s", err.Error()))
				}
			}

			if nextToken == "" {
				return nil // done -- TODO? tell the frontend the query is done
			}
			q.query.NextToken = nextToken
		}
	}
}
