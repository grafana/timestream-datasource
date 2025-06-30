package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/timestream-datasource/pkg/common"
)

// FormatQueryOption defines how the user has chosen to represent the data
type FormatQueryOption uint32

const (
	// FormatOptionTable formats the query results as a table using "LongToWide"
	FormatOptionTable FormatQueryOption = iota
	//FormatOptionTimeSeries formats the query results as a timeseries using "WideToLong"
	FormatOptionTimeSeries
)

var LegacyQueryCheck = regexp.MustCompile(`"format":\s*"table"`)

// QueryModel represents a spreadsheet query.
type QueryModel struct {
	RawQuery  string `json:"rawQuery,omitempty"`
	NextToken string `json:"nextToken,omitempty"`

	// Templates ${value}
	Database string `json:"database,omitempty"`
	Table    string `json:"table,omitempty"`
	Measure  string `json:"measure,omitempty"`

	// Not from JSON
	Interval      time.Duration     `json:"-"`
	TimeRange     backend.TimeRange `json:"-"`
	MaxDataPoints int64             `json:"-"`

	// Return several pages (if exist) in one response
	WaitForResult bool `json:"waitForResult"`

	// Format the results
	Format FormatQueryOption `json:"format"`
}

// GetQueryModel returns a parsed query
func GetQueryModel(query backend.DataQuery) (*QueryModel, error) {
	model := &QueryModel{}

	err := json.Unmarshal(query.JSON, &model)
	if err != nil {
		if LegacyQueryCheck.Match(query.JSON) {
			return nil, backend.DownstreamError(fmt.Errorf("query is incompatible with current structure, please rebuild it: %w", err))
		}
		return nil, backend.PluginError(fmt.Errorf("error reading query: %s", err.Error()))
	}

	// Copy directly from the well typed query
	model.TimeRange = query.TimeRange
	model.Interval = query.Interval
	model.MaxDataPoints = query.MaxDataPoints

	// In 7.1 alerting queries send empty values for MaxDataPoints
	if model.MaxDataPoints == 0 {
		model.MaxDataPoints = 1024
	}

	// In 7.1 alerting queries send empty values for interval
	if model.Interval.Milliseconds() == 0 && model.MaxDataPoints > 0 {
		millis := model.TimeRange.Duration().Milliseconds() / model.MaxDataPoints
		model.Interval = time.Millisecond * time.Duration(common.RoundInterval(millis))
	}

	return model, nil
}

// CancelRequest will cancel a running query
type CancelRequest struct {
	QueryID string `json:"queryId,omitempty"`
}

// TablesRequest will return tables for a database
type TablesRequest struct {
	Database string `json:"database"`
}

// CancelRequest will return measures for a table
type MeasuresRequest struct {
	Database string `json:"database"`
	Table    string `json:"table"`
}
