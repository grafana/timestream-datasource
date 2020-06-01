package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// QueryModel represents a spreadsheet query.
type QueryModel struct {
	RawQuery string   `json:"rawQuery"`
	Labels   []string `json:"labels"`

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
	model.Interval = query.Interval
	model.TimeRange = query.TimeRange
	model.MaxDataPoints = query.MaxDataPoints
	return model, nil
}
