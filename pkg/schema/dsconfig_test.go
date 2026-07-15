package schema_test

import (
	_ "embed"
	"testing"

	"github.com/grafana/dsconfig/schema"
	"github.com/grafana/timestream-datasource/pkg/models"
)

//go:embed dsconfig.json
var configSchemaJSON []byte

//go:generate go test -run TestPlugin -generateArtifacts
func TestPlugin(t *testing.T) {
	schema.RunPluginTests(t, schema.PluginUnderTest{
		ID:                "grafana-timestream-datasource",
		ConfigSchemaJSON:  configSchemaJSON,
		SettingsJSONModel: models.DatasourceSettings{},
		SecureKeys:        []string{"accessKey", "secretKey", "sessionToken"},
	})
}
