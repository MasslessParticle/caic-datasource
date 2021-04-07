package main

import (
	"os"

	"github.com/grafana/caic-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func main() {
	// This is the thing that is served
	ds := plugin.DatasourceOpts()
	err := datasource.Serve(ds)

	// og any error if we could start the plugin.
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
