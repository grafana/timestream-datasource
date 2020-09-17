package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/timestream-datasource/pkg/common"
)

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
}

// GetQueryModel returns a parsed query
func GetQueryModel(query backend.DataQuery) (*QueryModel, error) {
	model := &QueryModel{}

	err := json.Unmarshal(query.JSON, &model)
	if err != nil {
		return nil, fmt.Errorf("error reading query: %s", err.Error())
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
