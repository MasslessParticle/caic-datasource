package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grafana/caic-datasource/pkg/caic"
	"github.com/grafana/caic-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func main() {
	fmt.Println("Plugin Starting")
	caicURL := os.Getenv("CAIC_ADDR")
	fmt.Println("caic URL: ", caicURL)

	addr := os.Getenv("GF_PLUGIN_GRPC_ADDRESS_" + strings.ReplaceAll(strings.ToUpper("grafana-caic-datasource"), "-", "_"))
	fmt.Println("addr: ", addr)

	if err := datasource.Manage("grafana-caic-datasource", constructor, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

//TODO For this to work with the standalone stuff, the plugin needs
func constructor(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	caicURL := os.Getenv("CAIC_ADDR")
	if caicURL == "" {
		caicURL = "https://www.avalanche.state.co.us"
	}

	client := caic.NewClient(caicURL, http.DefaultClient)
	cache := caic.NewClientCache(client)
	return &plugin.Handler{
		Client: cache,
	}, nil
}
