package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamquery"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/timestream-datasource/pkg/models"
	"github.com/grafana/timestream-datasource/pkg/timestream"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"

	"context"
)

const metricNamespace = "timestream"

// TimestreamDataSource handler for google sheets
type TimestreamDataSource struct {
	Cache *cache.Cache
}

// NewDataSource creates the google sheets datasource and sets up all the routes
func NewDataSource(mux *http.ServeMux) *TimestreamDataSource {
	queriesTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "data_query_total",
			Help:      "data query counter",
			Namespace: metricNamespace,
		},
		[]string{"scenario"},
	)
	prometheus.MustRegister(queriesTotal)

	cache := cache.New(300*time.Second, 5*time.Second)
	ds := &TimestreamDataSource{
		Cache: cache,
	}

	mux.HandleFunc("/spreadsheets", ds.handleResourceSpreadsheets)
	return ds
}

func readConfig(pluginConfig backend.PluginConfig) (*models.TimestreamConfig, error) {
	config := models.TimestreamConfig{}
	if err := json.Unmarshal(pluginConfig.DataSourceConfig.JSONData, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal DataSourceInfo json: %w", err)
	}
	// config.APIKey = pluginConfig.DataSourceConfig.DecryptedSecureJSONData["apiKey"]
	// config.JWT = pluginConfig.DataSourceConfig.DecryptedSecureJSONData["jwt"]
	return &config, nil
}

func readQuery(q backend.DataQuery) (*models.QueryModel, error) {
	queryModel := models.QueryModel{}
	if err := json.Unmarshal(q.JSON, &queryModel); err != nil {
		return nil, fmt.Errorf("failed to unmarshal query: %w", err)
	}
	return &queryModel, nil
}

// GetSession for the current datasource
func (ds *TimestreamDataSource) GetSession(pluginConfig *models.TimestreamConfig) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	return sess, err
}

// CheckHealth checks if the plugin is running properly
func (ds *TimestreamDataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	res := &backend.CheckHealthResult{}

	// Just checking that the plugin exe is alive and running
	if req.PluginConfig.DataSourceConfig == nil {
		res.Status = backend.HealthStatusOk
		res.Message = "Plugin is Running"
		return res, nil
	}

	config, err := readConfig(req.PluginConfig)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Invalid config"
		return res, nil
	}

	sess, err := ds.GetSession(config)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Unable to get Session"
		return res, nil
	}

	if sess == nil {
		return nil, fmt.Errorf("null session")
	}

	res.Status = backend.HealthStatusOk
	res.Message = "Success"
	return res, nil
}

// QueryData queries for data.
func (ds *TimestreamDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	res := &backend.QueryDataResponse{}
	config, err := readConfig(req.PluginConfig)
	if err != nil {
		return nil, err
	}

	sess, err := ds.GetSession(config)
	if err != nil {
		return nil, err
	}

	for _, q := range req.Queries {
		queryModel, err := readQuery(q)
		if err != nil {
			return nil, fmt.Errorf("failed to read query: %w", err)
		}

		if len(queryModel.RawQuery) < 1 {
			continue // not query really exists
		}

		allowTruncation := !queryModel.NoTruncation
		queryInput := &timestreamquery.QueryInput{
			QueryString:           &queryModel.RawQuery,
			AllowResultTruncation: &allowTruncation,
		}
		backend.Logger.Debug("QueryInput:", queryModel.RawQuery)

		// execute the query
		querySvc := timestreamquery.New(sess)
		result, err := querySvc.Query(queryInput)
		if err != nil {
			// TODO: add a frame with an error header
			return nil, err
		}

		frame, err := timestream.QueryResultToDataFrame(result)
		if err != nil {
			// TODO: add a frame with an error header
			return nil, err
		}
		frame.Name = q.RefID // ???
		frame.RefID = q.RefID

		res.Frames = append(res.Frames, []*data.Frame{frame}...)
	}

	return res, nil
}

func writeResult(rw http.ResponseWriter, path string, val interface{}, err error) {
	response := make(map[string]interface{})
	code := http.StatusOK
	if err != nil {
		response["error"] = err.Error()
		code = http.StatusBadRequest
	} else {
		response[path] = val
	}

	body, err := json.Marshal(response)
	if err != nil {
		body = []byte(err.Error())
		code = http.StatusInternalServerError
	}
	_, err = rw.Write(body)
	if err != nil {
		code = http.StatusInternalServerError
	}
	rw.WriteHeader(code)
}

func (ds *TimestreamDataSource) handleResourceSpreadsheets(rw http.ResponseWriter, req *http.Request) {
	res := map[string]string{}
	res["hello"] = "world"
	writeResult(rw, "spreadsheets", res, nil)
}
