package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// DatasourceSettings holds basic connection info
type DatasourceSettings struct {
	awsds.AWSDatasourceSettings

	Config backend.DataSourceInstanceSettings

	// Default query
	DefaultDatabase string `json:"defaultDatabase,omitempty"`
	DefaultTable    string `json:"defaultTable,omitempty"`
	DefaultMeasure  string `json:"defaultMeasure,omitempty"`
}

// Load is copied from grafana-aws-sdk -- json.Unmarshal was not loading the nested properties
func (s *DatasourceSettings) Load(config backend.DataSourceInstanceSettings) error {
	s.Config = config
	if len(config.JSONData) > 1 {
		if err := json.Unmarshal(config.JSONData, s); err != nil {
			return fmt.Errorf("could not unmarshal DatasourceSettings json: %w", err)
		}
	}

	if s.Region == "default" || s.Region == "" {
		s.Region = s.DefaultRegion
	}

	if s.Profile == "" {
		s.Profile = config.Database // legacy support (only for cloudwatch?)
	}

	s.AccessKey = config.DecryptedSecureJSONData["accessKey"]
	s.SecretKey = config.DecryptedSecureJSONData["secretKey"]

	return nil
}
