package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/require"
)

func TestThePlugin(t *testing.T) {
	handler := &handler{}
	url, shutdown := startTestAPIServer(handler)
	defer shutdown()

	env := fmt.Sprintf("CAIC_ADDR=%s", url)
	plugin, cleanup := StartPlugin(t, "github.com/grafana/caic-datasource/pkg", []string{env})
	defer cleanup()

	t.Run("it returns success when the caic site is reachable", func(t *testing.T) {
		handler.status = http.StatusOK

		res, err := plugin.CheckHealth(context.Background(), healthReq)
		require.Nil(t, err)

		require.Equal(t, "Data source is working", res.Message)
	})

	t.Run("it returns an error when the caic site is unavailable", func(t *testing.T) {
		handler.status = http.StatusNotFound

		res, err := plugin.CheckHealth(context.Background(), healthReq)
		require.Nil(t, err)

		require.Equal(t, "Error reaching CAIC site", res.Message)
	})
}

func startTestAPIServer(h http.Handler) (string, func()) {
	server := httptest.NewServer(h)
	shutdown := func() {
		server.Config.Shutdown(context.Background())
	}
	return server.URL, shutdown
}

var healthReq = &backend.CheckHealthRequest{
	PluginContext: backend.PluginContext{
		DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
			ID: 0,
		},
	},
}

type handler struct {
	status int
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(h.status)
}
