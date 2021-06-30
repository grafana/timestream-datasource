package timestream

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/timestream-datasource/pkg/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
)

type clientGetterFunc func(region string) (client *timestreamquery.TimestreamQuery, err error)

// This is an interface to help testing
type queryRunner interface {
	runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error)
	cancelQuery(ctx context.Context, input *timestreamquery.CancelQueryInput) (*timestreamquery.CancelQueryOutput, error)
}

type timestreamRunner struct {
	querySvc clientGetterFunc
}

func (r *timestreamRunner) runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error) {
	svc, err := r.querySvc("")
	if err != nil {
		return nil, err
	}
	return svc.QueryWithContext(ctx, input)
}

func (r *timestreamRunner) cancelQuery(ctx context.Context, input *timestreamquery.CancelQueryInput) (*timestreamquery.CancelQueryOutput, error) {
	svc, err := r.querySvc("")
	if err != nil {
		return nil, err
	}
	return svc.CancelQueryWithContext(ctx, input)
}

func NewServerInstance(s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings := models.DatasourceSettings{}
	err := settings.Load(s)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}
	backend.Logger.Info("new instance", "settings", settings)
	sessions := awsds.NewSessionCache()
	endpoint := settings.Endpoint
	settings.Endpoint = "" // do not use this in the initial session configuration

	return &timestreamDS{
		Settings:      settings,
		streams:       make(map[string]*openQuery),
		channelPrefix: fmt.Sprintf("ds/%s/", s.UID),
		Runner: &timestreamRunner{
			querySvc: func(region string) (client *timestreamquery.TimestreamQuery, err error) {

				sess, err := sessions.GetSession(region, settings.AWSDatasourceSettings)
				if err != nil {
					return nil, err
				}
				tcfg := &aws.Config{}
				if endpoint != "" {
					tcfg.Endpoint = aws.String(endpoint)
				}
				querySvc := timestreamquery.New(sess, tcfg)

				// Add the user agent version
				querySvc.Handlers.Send.PushFront(func(r *request.Request) {
					r.HTTPRequest.Header.Set("User-Agent", fmt.Sprintf("GrafanaTimestream/%s", "1.0.0"))
				})
				return querySvc, nil
			},
		},
	}, nil
}

func (s *timestreamDS) Dispose() {
	// Called before creatinga a new instance to allow plugin authors
	// to cleanup.
}

type timestreamDS struct {
	Runner   queryRunner
	Settings models.DatasourceSettings

	mu            sync.RWMutex
	streams       map[string]*openQuery
	channelPrefix string
}

var (
	_ backend.QueryDataHandler      = (*timestreamDS)(nil)
	_ backend.CheckHealthHandler    = (*timestreamDS)(nil)
	_ backend.StreamHandler         = (*timestreamDS)(nil)
	_ instancemgmt.InstanceDisposer = (*timestreamDS)(nil)
)

// CheckHealth will check the currently configured settings
func (ds *timestreamDS) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	// Connection is OK
	input := &timestreamquery.QueryInput{
		QueryString: aws.String("SELECT 1"),
	}
	output, err := ds.Runner.runQuery(ctx, input)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	val := output.Rows[0].Data[0].ScalarValue
	if val == nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "missing response",
		}, nil
	}
	if *val != "1" {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "should be one",
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Connection success",
	}, nil
}

// QueryData - Primary method called by grafana-server
func (ds *timestreamDS) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	res := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		query, err := models.GetQueryModel(q)
		if err != nil {
			res.Responses[q.RefID] = backend.DataResponse{
				Error: err,
			}
		} else {
			dr, nextToken := ExecuteQuery(ctx, *query, ds.Runner, ds.Settings)

			// If the query has a nextToken, register a stream to keep check pages
			if nextToken != "" {
				firstFrame := dr.Frames[0]
				firstFrame.Meta.Channel = ds.registerPagingQuery(firstFrame, query)
			}
			res.Responses[q.RefID] = dr
		}
	}
	return res, nil
}

// CallResource HTTP style resource
func (ds *timestreamDS) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	if req.Path == "hello" {
		return resource.SendPlainText(sender, "world")
	}
	if req.Path == "cancel" {
		if req.Method != "POST" {
			return fmt.Errorf("Cancel requires a post command")
		}
		cancel := &models.CancelRequest{}

		err := json.Unmarshal(req.Body, &cancel)
		if err != nil {
			return fmt.Errorf("error reading cancel request: %s", err.Error())
		}
		cancelQueryInput := &timestreamquery.CancelQueryInput{
			QueryId: aws.String(cancel.QueryID),
		}
		msg := "cancel: " + cancel.QueryID
		v, err := ds.Runner.cancelQuery(ctx, cancelQueryInput)
		if v != nil && v.CancellationMessage != nil {
			msg = *v.CancellationMessage
		} else if err != nil {
			msg = err.Error()
		}
		return resource.SendPlainText(sender, msg)
	}
	return fmt.Errorf("unknown resource")
}

func (ds *timestreamDS) registerPagingQuery(frame *data.Frame, query *models.QueryModel) string {
	meta := frame.Meta.Custom.(*models.TimestreamCustomMeta)
	if meta.NextToken == "" {
		return "" // no channel
	}
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// add stream query
	query.NextToken = meta.NextToken
	ds.streams[meta.QueryID] = &openQuery{
		queryId: meta.QueryID,
		query:   query,
		ds:      ds,
	}
	return ds.channelPrefix + meta.QueryID
}

func (ds *timestreamDS) SubscribeStream(_ context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	rsp := &backend.SubscribeStreamResponse{
		Status: backend.SubscribeStreamStatusNotFound,
	}

	ds.mu.RLock()
	defer ds.mu.RUnlock()

	s, ok := ds.streams[req.Path]
	if s != nil && ok {
		rsp.Status = backend.SubscribeStreamStatusOK
	}
	return rsp, nil
}

func (ds *timestreamDS) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	ds.mu.RLock()
	s, ok := ds.streams[req.Path]
	if s == nil || !ok {
		ds.mu.RUnlock()
		// Return nil to stop RunStream till next subscriber. Any error here
		// will result into RunStream re-establishment.
		return nil
	}
	ds.mu.RUnlock()

	// When the stream is done, remove it.
	defer func() {
		ds.mu.Lock()
		defer ds.mu.Unlock()
		delete(ds.streams, req.Path)
	}()

	return s.doStream(ctx, sender)
}

func (ds *timestreamDS) PublishStream(_ context.Context, _ *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	return &backend.PublishStreamResponse{
		Status: backend.PublishStreamStatusPermissionDenied,
	}, nil
}
