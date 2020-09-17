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
