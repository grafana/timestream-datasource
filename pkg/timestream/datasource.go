package timestream

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource"
	"github.com/grafana/timestream-datasource/pkg/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
)

// This is an interface to help testing
type queryRunner interface {
	runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error)
}

// This is an interface to help testing
type TimestreamRunner struct {
	querySvc *timestreamquery.TimestreamQuery
}

func (r *TimestreamRunner) runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error) {
	return r.querySvc.Query(input)
}

type instanceSettings struct {
	Runner queryRunner
}

func newDataSourceInstance(s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings, err := models.LoadSettings(s)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}
	backend.Logger.Info("new instance", "settings", settings)

	// setup the query client
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		return nil, err
	}

	// querySvc := timestreamquery.New(sess, &aws.Config{Endpoint: aws.String("query-cell0.timestream.us-east-1.amazonaws.com")})
	querySvc := timestreamquery.New(sess)

	return &instanceSettings{
		Runner: &TimestreamRunner{
			querySvc: querySvc,
		},
	}, nil
}

func (s *instanceSettings) Dispose() {
	// Called before creatinga a new instance to allow plugin authors
	// to cleanup.
}

// TimestreamDS ...
type TimestreamDS struct {
	im instancemgmt.InstanceManager
}

// NewDatasource creates a new datasource server
func NewDatasource() datasource.ServeOpts {
	im := datasource.NewInstanceManager(newDataSourceInstance)
	ds := &TimestreamDS{
		im: im,
	}

	return datasource.ServeOpts{
		QueryDataHandler:    ds,
		CheckHealthHandler:  ds,
		CallResourceHandler: ds,
	}
}

func (ds *TimestreamDS) getInstance(ctx backend.PluginContext) (*instanceSettings, error) {
	s, err := ds.im.Get(ctx)
	if err != nil {
		return nil, err
	}
	return s.(*instanceSettings), nil // ugly cast... but go ¯\_(ツ)_/¯
}

// CheckHealth will check the currently configured settings
func (ds *TimestreamDS) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	s, err := ds.getInstance(req.PluginContext)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	backend.Logger.Info("check health", "settings", s)

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: fmt.Sprintf("OK!"),
	}, nil
}

// QueryData - Primary method called by grafana-server
func (ds *TimestreamDS) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
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
			res.Responses[q.RefID] = ExecuteQuery(context.Background(), *query, s.Runner)
		}
	}
	return res, nil
}

// CallResource HTTP style resource
func (ds *TimestreamDS) CallResource(tx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	if req.Path == "hello" {
		return resource.SendPlainText(sender, "world")
	}

	return fmt.Errorf("unknown resource")
}
