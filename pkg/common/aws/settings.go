package aws

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// DatasourceSettings holds basic connection info
type DatasourceSettings struct {
	Profile       string `json:"profile"`
	Region        string `json:"region"`
	DefaultRegion string `json:"defaultRegion"`
	AuthType      string `json:"authType"`
	AssumeRoleArn string `json:"assumeRoleArn"`
	Namespace     string `json:"namespace"`

	// Timestream
	Endpoint string `json:"endpoint"`

	// Default query
	DefaultDatabase string `json:"defaultDatabase,omitempty"`
	DefaultTable    string `json:"defaultTable,omitempty"`
	DefaultMeasure  string `json:"defaultMeasure,omitempty"`

	// Loaded from DecryptedSecureJSONData (not the json object)
	AccessKey string `json:"-"`
	SecretKey string `json:"-"`
}

// LoadSettings will read and validate Settings from the DataSourceConfg
func LoadSettings(config backend.DataSourceInstanceSettings) (DatasourceSettings, error) {
	settings := DatasourceSettings{}

	if config.JSONData != nil && len(config.JSONData) > 1 {
		if err := json.Unmarshal(config.JSONData, &settings); err != nil {
			return settings, fmt.Errorf("could not unmarshal DatasourceSettings json: %w", err)
		}
	}

	if settings.Region == "default" || settings.Region == "" {
		settings.Region = settings.DefaultRegion
	}

	settings.AccessKey = config.DecryptedSecureJSONData["accessKey"]
	settings.SecretKey = config.DecryptedSecureJSONData["secretKey"]

	return settings, nil
}
