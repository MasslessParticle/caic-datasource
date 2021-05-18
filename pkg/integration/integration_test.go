package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/plugintest"
	"github.com/stretchr/testify/require"
)

func TestThePlugin(t *testing.T) {
	handler := &handler{}
	url, shutdown := startTestAPIServer(handler)
	defer shutdown()

	env := fmt.Sprintf("CAIC_ADDR=%s", url)
	client, cleanup, err := plugintest.StartPlugin(
		"github.com/grafana/caic-datasource/pkg",
		"grafana-caic-datasource",
		8000,
		env,
	)
	require.Nil(t, err)

	defer cleanup()

	t.Run("it returns success when the caic site is reachable", func(t *testing.T) {
		handler.status = http.StatusOK

		res, err := client.CheckHealth(context.Background(), healthReq)
		require.Nil(t, err)

		require.Equal(t, "Data source is working", res.Message)
	})

	t.Run("it returns an error when the caic site is unavailable", func(t *testing.T) {
		handler.status = http.StatusNotFound

		res, err := client.CheckHealth(context.Background(), healthReq)
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
		PluginID: "grafana-caic-datasource",
		DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
			ID: 10,
		},
	},
}

type handler struct {
	status int
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(h.status)
}
