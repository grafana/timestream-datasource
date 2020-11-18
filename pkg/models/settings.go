package models

import (
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
)

// DatasourceSettings holds basic connection info
type DatasourceSettings struct {
	awsds.AWSDatasourceSettings

	// Default query
	DefaultDatabase string `json:"defaultDatabase,omitempty"`
	DefaultTable    string `json:"defaultTable,omitempty"`
	DefaultMeasure  string `json:"defaultMeasure,omitempty"`
}
