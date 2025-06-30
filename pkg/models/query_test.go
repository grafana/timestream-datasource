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
