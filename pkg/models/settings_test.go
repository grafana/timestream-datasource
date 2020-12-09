package models

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestReadSettings(t *testing.T) {

	s := backend.DataSourceInstanceSettings{
		ID: 33,
		JSONData: []byte(`{
			"authType": "keys",
			"defaultDatabase": "sampleDB",
			"defaultMeasure": "speed",
			"defaultRegion": "us-west-2",
			"defaultTable": "IoT"
		  }`),
	}

	settings := DatasourceSettings{}
	err := settings.Load(s)
	if err != nil {
		t.Fatal("should not error")
	}

	if settings.DefaultDatabase != "sampleDB" {
		t.Fatalf("invalid data points: %s", settings.DefaultDatabase)
	}
}
