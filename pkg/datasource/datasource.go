package datasource

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/timestream-datasource/pkg/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
)

// This is an interface to help testing
type queryRunner interface {
	runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error)
}

// TimestreamDataSource handler for google sheets
type TimestreamDataSource struct {
	Runner queryRunner
}

// This is an interface to help testing
type TimestreamRunner struct {
	querySvc *timestreamquery.TimestreamQuery
}

func (r *TimestreamRunner) runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error) {
	return r.querySvc.Query(input)
}

// CreateDataSource create the client...
func CreateDataSource(settings *models.DatasourceSettings) (*TimestreamDataSource, error) {

	// setup the query client
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		return nil, err
	}

	// querySvc := timestreamquery.New(sess, &aws.Config{Endpoint: aws.String("query-cell0.timestream.us-east-1.amazonaws.com")})
	querySvc := timestreamquery.New(sess)

	return &TimestreamDataSource{
		Runner: &TimestreamRunner{
			querySvc: querySvc,
		},
	}, nil
}

// CheckHealth will check the currently configured settings
func (ds *TimestreamDataSource) CheckHealth() *backend.CheckHealthResult {

	// return &backend.CheckHealthResult{
	// 	Status:  backend.HealthStatusError,
	// 	Message: err.Error(),
	// }

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: fmt.Sprintf("OK!"),
	}
}

// QueryData - Primary method called by grafana-server
func (ds *TimestreamDataSource) QueryData(req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	res := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		query, err := models.GetQueryModel(q)
		if err != nil {
			res.Responses[q.RefID] = backend.DataResponse{
				Error: err,
			}
		} else {
			res.Responses[q.RefID] = ExecuteQuery(context.Background(), *query, ds.Runner)
		}
	}
	return res, nil
}

// CallResource HTTP style resource
func (ds *TimestreamDataSource) CallResource(req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {

	if req.Path == "hello" {
		return experimental.SendPlainText(sender, "world")
	}

	return fmt.Errorf("unknown resource")
}

// Destroy destroy an instance (if necessary)
func (ds *TimestreamDataSource) Destroy() {
	// If necessary, destroy the object (typically not required)
}
