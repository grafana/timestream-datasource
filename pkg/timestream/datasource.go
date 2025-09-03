package timestream

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana-aws-sdk/pkg/awsauth"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/errorsource"
	"github.com/grafana/timestream-datasource/pkg/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
	timestreamquerytypes "github.com/aws/aws-sdk-go-v2/service/timestreamquery/types"
)

type QueryClient interface {
	timestreamquery.QueryAPIClient
	CancelQuery(context.Context, *timestreamquery.CancelQueryInput, ...func(*timestreamquery.Options)) (*timestreamquery.CancelQueryOutput, error)
}

func NewDatasource(ctx context.Context, s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings := models.DatasourceSettings{}
	err := settings.Load(s)
	if err != nil {
		return nil, errorsource.PluginError(fmt.Errorf("error reading settings: %s", err.Error()), false)
	}

	httpClientProvider := sdkhttpclient.NewProvider()
	httpClientOptions, err := settings.Config.HTTPClientOptions(ctx)
	if err != nil {
		backend.Logger.Error("failed to create HTTP client options", "error", err.Error())
		return nil, errorsource.PluginError(err, false)
	}
	httpClient, err := httpClientProvider.New(httpClientOptions)
	if err != nil {
		backend.Logger.Error("failed to create HTTP client", "error", err.Error())
		return nil, errorsource.PluginError(err, false)
	}
	region := settings.Region
	if region == "" || region == "default" {
		region = settings.DefaultRegion
	}
	cfg, err := awsauth.NewConfigProvider().GetConfig(ctx, awsauth.Settings{
		LegacyAuthType:     settings.AuthType,
		AccessKey:          settings.AccessKey,
		SecretKey:          settings.SecretKey,
		Region:             region,
		CredentialsProfile: settings.Profile,
		AssumeRoleARN:      settings.AssumeRoleARN,
		Endpoint:           settings.Endpoint,
		ExternalID:         settings.ExternalID,
		UserAgent:          "Timestream",
		HTTPClient:         httpClient,
	})
	if err != nil {
		return nil, backend.DownstreamError(err)
	}

	var client QueryClient
	if settings.Endpoint != "" && settings.Endpoint != "default" {
		client = timestreamquery.NewFromConfig(cfg, func(o *timestreamquery.Options) {
			// Why disable Endpoint Discovery when a custom endpoint (e.g., VPC endpoint) is configured?
			// - AWS SDK for Go v1: endpoint discovery was OFF by default and only used when endpoint was unset/empty.
			//   VPC setups worked because DescribeEndpoints was not invoked.
			// - AWS SDK for Go v2: endpoint discovery defaults to AUTO. Because Timestream requires discovery,
			//   the SDK will call DescribeEndpoints before Query. With a custom BaseEndpoint (VPC endpoint),
			//   DescribeEndpoints is routed through that endpoint, which typically does not implement it,
			//   resulting in a 404 response.
			// - Even forcing discovery through the SDK’s default public resolver can fail if the VPC blocks
			//   egress to public AWS endpoints.
			// To preserve existing customer VPC configurations and avoid breaking changes, we explicitly disable
			// endpoint discovery whenever a custom endpoint is provided. Regular operations still use the custom endpoint.
			//
			// References:
			// - v1 behavior (see `EnableEndpointDiscovery` default): https://docs.aws.amazon.com/sdk-for-go/api/aws/
			// - v2 EndpointDiscoveryEnableState (AUTO default): https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws#EndpointDiscoveryEnableState
			o.EndpointDiscovery.EnableEndpointDiscovery = aws.EndpointDiscoveryDisabled
		})
	} else {
		client = timestreamquery.NewFromConfig(cfg)
	}

	return &timestreamDS{
		Settings: settings,
		Client:   client,
	}, nil
}

type timestreamDS struct {
	Client   QueryClient
	Settings models.DatasourceSettings
}

var (
	_ backend.QueryDataHandler   = (*timestreamDS)(nil)
	_ backend.CheckHealthHandler = (*timestreamDS)(nil)
)

// CheckHealth will check the currently configured settings
func (ds *timestreamDS) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	// Connection is OK
	input := &timestreamquery.QueryInput{
		QueryString: aws.String("SELECT 1"),
	}
	output, err := ds.Client.Query(ctx, input)
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
			errorsource.AddErrorToResponse(q.RefID, res, err)
		} else {
			res.Responses[q.RefID] = ds.ExecuteQuery(ctx, *query)
		}
	}
	return res, nil
}

func sliceFromRows(rows []timestreamquerytypes.Row, doubleQuotes bool) []string {
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
func dimensionsFromRows(rows []timestreamquerytypes.Row) []string {
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
			return fmt.Errorf("cancel requires a post command")
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
		v, err := ds.Client.CancelQuery(ctx, cancelQueryInput)
		if v != nil && v.CancellationMessage != nil {
			msg = *v.CancellationMessage
		} else if err != nil {
			msg = err.Error()
		}
		return resource.SendPlainText(sender, msg)
	}
	if req.Path == "databases" {
		// TODO: Use API endpoint to list databases
		v, err := ds.Client.Query(ctx, &timestreamquery.QueryInput{
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
			return fmt.Errorf("tables requires a post command")
		}
		opts := models.TablesRequest{}
		err := json.Unmarshal(req.Body, &opts)
		if err != nil {
			return err
		}
		// TODO: Use API endpoint to list tables
		v, err := ds.Client.Query(ctx, &timestreamquery.QueryInput{
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
			return fmt.Errorf("measures requires a post command")
		}
		opts := models.MeasuresRequest{}
		err := json.Unmarshal(req.Body, &opts)
		if err != nil {
			return err
		}
		v, err := ds.Client.Query(ctx, &timestreamquery.QueryInput{
			QueryString: aws.String(fmt.Sprintf("SHOW MEASURES FROM %s.%s", applyQuotesIfNeeded(opts.Database), applyQuotesIfNeeded(opts.Table))),
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

func applyQuotesIfNeeded(input string) string {
	if input[0] != '"' && input[len(input)-1] != '"' {
		input = fmt.Sprintf(`"%s"`, input)
	}
	return input
}

// ExecuteQuery -- run a query
func (ds *timestreamDS) ExecuteQuery(ctx context.Context, query models.QueryModel) backend.DataResponse {
	raw, err := Interpolate(query, ds.Settings)
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
	output, err := ds.Client.Query(ctx, input)
	if err == nil && query.WaitForResult && output.NextToken != nil {
		for output.NextToken != nil {
			newPageInput := *input
			newPageInput.NextToken = output.NextToken
			newPageOutput, newPageErr := ds.Client.Query(ctx, &newPageInput)
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
