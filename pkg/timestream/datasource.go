package timestream

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource"
	"github.com/grafana/timestream-datasource/pkg/models"

	"github.com/aws/aws-sdk-go/aws"
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
	sessions := awsds.NewSessionCache()
	endpoint := settings.Endpoint
	settings.Endpoint = "" // do not use this in the initial session configuration

	return &timestreamDS{
		Settings: settings,
		Runner: &timestreamRunner{
			querySvc: func(region string) (client *timestreamquery.TimestreamQuery, err error) {

				httpClientProvider := sdkhttpclient.NewProvider()
				httpClientOptions, err := settings.Config.HTTPClientOptions()
				if err != nil {
					backend.Logger.Error("failed to create HTTP client options", "error", err.Error())
					return nil, err
				}
				httpClient, err := httpClientProvider.New(httpClientOptions)
				if err != nil {
					backend.Logger.Error("failed to create HTTP client", "error", err.Error())
					return nil, err
				}

				sess, err := sessions.GetSession(awsds.SessionConfig{
					Settings:      settings.AWSDatasourceSettings,
					HTTPClient:    httpClient,
					UserAgentName: aws.String("Timestream"),
				})
				if err != nil {
					return nil, err
				}
				tcfg := &aws.Config{}
				if endpoint != "" {
					tcfg.Endpoint = aws.String(endpoint)
				}
				querySvc := timestreamquery.New(sess, tcfg)

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
}

var (
	_ backend.QueryDataHandler      = (*timestreamDS)(nil)
	_ backend.CheckHealthHandler    = (*timestreamDS)(nil)
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
			res.Responses[q.RefID] = ExecuteQuery(ctx, *query, ds.Runner, ds.Settings)
		}
	}
	return res, nil
}

func sliceFromRows(rows []*timestreamquery.Row, doubleQuotes bool) []string {
	res := []string{}
	for _, row := range rows {
		if len(row.Data) > 0 && row.Data[0].ScalarValue != nil {
			val := *row.Data[0].ScalarValue
			if doubleQuotes {
				val = fmt.Sprintf(`"%s"`, val)
			}
			res = append(res, val)
		}
	}
	return res
}

// The dimensions of a row are the different columns (apart from the measure)
// It's encoded in every row, as an array of values.
func dimensionsFromRows(rows []*timestreamquery.Row) []string {
	res := []string{}
	if len(rows) > 0 && len(rows[0].Data) == 3 && len(rows[0].Data[2].ArrayValue) > 0 {
		dimensionArray := rows[0].Data[2].ArrayValue
		for _, dim := range dimensionArray {
			if len(dim.RowValue.Data) > 0 && dim.RowValue.Data[0].ScalarValue != nil {
				res = append(res, *dim.RowValue.Data[0].ScalarValue)
			}
		}
	}
	return res
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
	if req.Path == "databases" {
		// TODO: Use API endpoint to list databases
		v, err := ds.Runner.runQuery(ctx, &timestreamquery.QueryInput{
			QueryString: aws.String("SHOW DATABASES"),
		})
		if err != nil {
			return err
		}
		// Databases are returned wrapped in double quotes
		return resource.SendJSON(sender, sliceFromRows(v.Rows, true))
	}
	if req.Path == "tables" {
		if req.Method != "POST" {
			return fmt.Errorf("Tables requires a post command")
		}
		opts := models.TablesRequest{}
		err := json.Unmarshal(req.Body, &opts)
		if err != nil {
			return err
		}
		// TODO: Use API endpoint to list tables
		v, err := ds.Runner.runQuery(ctx, &timestreamquery.QueryInput{
			QueryString: aws.String(fmt.Sprintf("SHOW TABLES FROM %s", applyQuotesIfNeeded(opts.Database))),
		})
		if err != nil {
			return err
		}
		// Tables are returned wrapped in double quotes
		return resource.SendJSON(sender, sliceFromRows(v.Rows, true))
	}
	if req.Path == "measures" || req.Path == "dimensions" {
		if req.Method != "POST" {
			return fmt.Errorf("Measures requires a post command")
		}
		opts := models.MeasuresRequest{}
		err := json.Unmarshal(req.Body, &opts)
		if err != nil {
			return err
		}
		v, err := ds.Runner.runQuery(ctx, &timestreamquery.QueryInput{
			QueryString: aws.String(fmt.Sprintf("SHOW MEASURES FROM %s.%s", opts.Database, opts.Table)),
		})
		if err != nil {
			return err
		}
		if req.Path == "measures" {
			return resource.SendJSON(sender, sliceFromRows(v.Rows, false))
		}
		if req.Path == "dimensions" {
			return resource.SendJSON(sender, dimensionsFromRows(v.Rows))
		}
	}
	return fmt.Errorf("unknown resource")
}

func applyQuotesIfNeeded(dbName string) string {
	if dbName[0] != '"' && dbName[len(dbName)-1] != '"' {
		dbName = fmt.Sprintf(`"%s"`, dbName)
	}
	return dbName
}
