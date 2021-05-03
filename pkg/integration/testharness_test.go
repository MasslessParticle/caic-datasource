package integration_test

import (
	"context"
	"testing"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana/pkg/plugins/backendplugin/grpcplugin"
	"github.com/onsi/gomega/gexec"
)

type Plugin interface {
	PluginID() string
	Logger() log.Logger
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	IsManaged() bool
	Exited() bool
	CollectMetrics(ctx context.Context) (*backend.CollectMetricsResult, error)
	CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error)
	CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error
	SubscribeStream(ctx context.Context, request *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error)
	PublishStream(ctx context.Context, request *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error)
	RunStream(ctx context.Context, req *backend.RunStreamRequest, sender backend.StreamPacketSender) error
}

var startCallbacks = grpcplugin.PluginStartFuncs{
	OnLegacyStart: func(string, *grpcplugin.LegacyClient, log.Logger) error { return nil },
	OnStart:       func(string, *grpcplugin.Client, log.Logger) error { return nil },
}

func StartPlugin(t *testing.T, packagePath string, env []string) (Plugin, func()) {
	pluginPath, err := gexec.Build(packagePath)
	require.Nil(t, err)

	pf := grpcplugin.NewBackendPlugin("test-plugin", pluginPath, startCallbacks)
	plugin, err := pf("test-plugin", log.New("Plugin Logger"), env)
	require.Nil(t, err)

	err = plugin.Start(context.Background())
	require.Nil(t, err)

	cleanup := func() {
		plugin.Stop(context.Background())
		gexec.CleanupBuildArtifacts()
	}

	return plugin, cleanup
}
