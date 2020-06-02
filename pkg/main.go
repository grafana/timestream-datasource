package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/timestream-datasource/pkg/timestream"
)

func main() {
	backend.SetupPluginEnvironment("timestream-datasource")

	err := datasource.Serve(timestream.NewDatasource())

	// Log any error if we could start the plugin.
	if err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
