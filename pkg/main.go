package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/timestream-datasource/pkg/timestream"
)

func main() {

	if err := datasource.Manage("timestream-datasource", timestream.NewServerInstance, datasource.ManageOpts{}); err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
