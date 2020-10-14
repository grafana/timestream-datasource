package timestream

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource"
	gaws "github.com/grafana/timestream-datasource/pkg/common/aws"
	"github.com/grafana/timestream-datasource/pkg/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
)

type instanceSettings struct {
	Runner   queryRunner
	Settings gaws.DatasourceSettings
}

// This is an interface to help testing
type queryRunner interface {
	runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error)
	cancelQuery(ctx context.Context, input *timestreamquery.CancelQueryInput) (*timestreamquery.CancelQueryOutput, error)
}

type timestreamRunner struct {
	querySvc *timestreamquery.TimestreamQuery
}

func (r *timestreamRunner) runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error) {
	return r.querySvc.QueryWithContext(ctx, input)
}

func (r *timestreamRunner) cancelQuery(ctx context.Context, input *timestreamquery.CancelQueryInput) (*timestreamquery.CancelQueryOutput, error) {
	return r.querySvc.CancelQueryWithContext(ctx, input)
}

func newDataSourceInstance(s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings, err := gaws.LoadSettings(s)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}
	backend.Logger.Info("new instance", "settings", settings)

	// setup the query client
	cfg, err := gaws.GetAwsConfig(settings)
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	tcfg := &aws.Config{}
	if settings.Endpoint != "" {
		tcfg.Endpoint = aws.String(settings.Endpoint)
	}
	querySvc := timestreamquery.New(sess, tcfg)

	// Add the user agent version
	querySvc.Handlers.Send.PushFront(func(r *request.Request) {
		r.HTTPRequest.Header.Set("User-Agent", fmt.Sprintf("GrafanaTimestream/%s", "alpha"))
	})

	return &instanceSettings{
		Settings: settings,
		Runner: &timestreamRunner{
			querySvc: querySvc,
		},
	}, nil
}

func (s *instanceSettings) Dispose() {
	// Called before creatinga a new instance to allow plugin authors
	// to cleanup.
}

type timestreamDS struct {
	im instancemgmt.InstanceManager
}

// NewDatasource creates a new datasource server
func NewDatasource() datasource.ServeOpts {
	im := datasource.NewInstanceManager(newDataSourceInstance)
	ds := &timestreamDS{
		im: im,
	}

	return datasource.ServeOpts{
		QueryDataHandler:    ds,
		CheckHealthHandler:  ds,
		CallResourceHandler: ds,
	}
}

func (ds *timestreamDS) getInstance(ctx backend.PluginContext) (*instanceSettings, error) {
	s, err := ds.im.Get(ctx)
	if err != nil {
		return nil, err
	}
	return s.(*instanceSettings), nil
}

// CheckHealth will check the currently configured settings
func (ds *timestreamDS) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	s, err := ds.getInstance(req.PluginContext)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	// Connection is OK
	input := &timestreamquery.QueryInput{
		QueryString: aws.String("SELECT 1"),
	}
	output, err := s.Runner.runQuery(ctx, input)
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
	s, err := ds.getInstance(req.PluginContext)
	if err != nil {
		return nil, err
	}

	res := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		query, err := models.GetQueryModel(q)
		if err != nil {
			res.Responses[q.RefID] = backend.DataResponse{
				Error: err,
			}
		} else {
			res.Responses[q.RefID] = ExecuteQuery(ctx, *query, s.Runner, s.Settings)
		}
	}
	return res, nil
}

// CallResource HTTP style resource
func (ds *timestreamDS) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	s, err := ds.getInstance(req.PluginContext)
	if err != nil {
		return err
	}

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
		v, err := s.Runner.cancelQuery(ctx, cancelQueryInput)
		if v != nil && v.CancellationMessage != nil {
			msg = *v.CancellationMessage
		} else if err != nil {
			msg = err.Error()
		}
		return resource.SendPlainText(sender, msg)
	}
	return fmt.Errorf("unknown resource")
}
