package main

import (
	"net/http"
	"os"

	"github.com/grafana/caic-datasource/pkg/caic"
	"github.com/grafana/caic-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func main() {
	im := datasource.NewInstanceManager(constructor)
	ds := plugin.DatasourceOpts(im)
	err := datasource.Serve(ds)

	// og any error if we could start the plugin.
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

func constructor(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	caicURL := "https://www.avalanche.state.co.us"
	return &plugin.CaicDatasource{
		Client: caic.NewClient(caicURL, http.DefaultClient),
	}, nil
}
