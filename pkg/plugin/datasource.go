package plugin

import (
	"net/http"

	"github.com/grafana/caic-datasource/pkg/caic"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

type CaicDatasource struct {
	client *caic.Client
}

func caicDataSourceInstance(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	caicURL := "https://www.avalanche.state.co.us"
	return &CaicDatasource{
		client: caic.NewClient(caicURL, http.DefaultClient),
	}, nil
}

func (s *CaicDatasource) Dispose() {
	// Called before creatinga a new instance to allow plugin authors
	// to cleanup.
}
