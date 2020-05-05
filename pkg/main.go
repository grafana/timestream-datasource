package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
)

func main() {
	// Setup the plugin environment
	backend.SetupPluginEnvironment("timestream-datasource")

	host := experimental.NewInstanceManager(&TimestreamHost{})
	err := host.RunGRPCServer()
	if err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
