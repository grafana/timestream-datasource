package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// InfluxSettings contains config properties (share with other AWS services?)
type DatasourceSettings struct {
}

// // Whether to use GZip compression in requests. Default false
// useGZip bool
// // TLS configuration for secure connection. Default nil
// tlsConfig *tls.Config
// // HTTP request timeout in sec. Default 20
// httpRequestTimeout uint

func LoadSettings(settings backend.DataSourceInstanceSettings) (*DatasourceSettings, error) {
	model := &DatasourceSettings{}

	err := json.Unmarshal(settings.JSONData, &model)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}

	return model, nil
}
