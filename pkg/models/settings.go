package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// Settings - all filled from
type Settings struct {
	Region string `json:"region"`
}

// LoadSettings will read and validate Settings from the DataSourceConfg
func LoadSettings(config backend.DataSourceInstanceSettings) (Settings, error) {
	settings := Settings{}

	if err := json.Unmarshal(config.JSONData, &settings); err != nil {
		return settings, fmt.Errorf("could not unmarshal DataSourceInfo json: %w", err)
	}

	return settings, nil
}
