package models

import (
	"encoding/json"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestAutomaticInterval(t *testing.T) {
	model := &QueryModel{}
	json, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("Unable to write json: %s", err.Error())
	}

	model, err = GetQueryModel(backend.DataQuery{
		JSON: json,
	})
	if err != nil {
		t.Fatalf("Error reading query: %s", err.Error())
	}

	if model.MaxDataPoints != 1024 {
		t.Fatalf("invalid data points: %d", model.MaxDataPoints)
	}
	if model.Interval.Milliseconds() != 10 {
		t.Fatalf("invalid interval: %d", model.Interval.Milliseconds())
	}
}

func TestRefIDIsCopiedFromQuery(t *testing.T) {
	model, err := GetQueryModel(backend.DataQuery{
		RefID: "First Sensor",
		JSON:  []byte(`{"rawQuery": "SELECT measure_value::double AS $__refId"}`),
	})
	if err != nil {
		t.Fatalf("Error reading query: %s", err.Error())
	}

	if model.RefID != "First Sensor" {
		t.Fatalf("invalid refId: %s", model.RefID)
	}
}

func TestGetQueryModel_Errors(t *testing.T) {
	tests := []struct {
		name           string
		rawQuery       string
		wantDownstream bool
	}{
		{
			name:           "empty query is plugin error",
			rawQuery:       "",
			wantDownstream: false,
		},
		{
			name:           "",
			rawQuery:       `{"format": "table", "group": [], "intervalMs": 1000, "maxDataPoints": 43200, "metricColumn": "none", "rawQuery": true, "rawSql": "select 1", "refId": "C", "select": [[{"params": ["id"], "type": "column"}]], "table": "a_table", "timeColumn": "auto_farmer_timestamp", "timeColumnType": "timestamp", "where": [{"name": "$__timeFilter", "params": [], "type": "macro"}]}`,
			wantDownstream: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := backend.DataQuery{
				JSON: []byte(tt.rawQuery),
			}
			_, err := GetQueryModel(query)
			if err == nil {
				t.Errorf("GetQueryModel() should have errored")
				return
			}
			isDownstream := backend.IsDownstreamError(err)
			if isDownstream != tt.wantDownstream {
				t.Errorf("GetQueryModel() error source isDownstream = %v, should be %v", isDownstream, tt.wantDownstream)
			}
		})
	}
}

func TestGetQueryModel_Format(t *testing.T) {
	tests := []struct {
		name string
		json string
		want FormatQueryOption
	}{
		{"missing", `{"rawQuery": "SELECT 1"}`, FormatOptionTable},
		{"numeric table", `{"rawQuery": "SELECT 1", "format": 0}`, FormatOptionTable},
		{"numeric time series", `{"rawQuery": "SELECT 1", "format": 1}`, FormatOptionTimeSeries},
		{"legacy time_series string", `{"rawQuery": "SELECT 1", "format": "time_series"}`, FormatOptionTable},
		{"legacy table string", `{"rawQuery": "SELECT 1", "format": "table"}`, FormatOptionTable},
		{"null", `{"rawQuery": "SELECT 1", "format": null}`, FormatOptionTable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := GetQueryModel(backend.DataQuery{JSON: []byte(tt.json)})
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if model.Format != tt.want {
				t.Fatalf("format = %d, want %d", model.Format, tt.want)
			}
		})
	}
}
