package plugin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/grafana/caic-datasource/pkg/caic"
	"github.com/grafana/caic-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/stretchr/testify/require"
)

func TestQueryData(t *testing.T) {
	t.Run("returns all zones when zone is -1", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.zones = []caic.Zone{
			{ID: "zone-1", Name: "zone 1", Rating: 1},
			{ID: "zone-2", Name: "zone 2", Rating: 3},
		}

		opts := plugin.DatasourceOpts(im)
		res, _ := opts.QueryDataHandler.QueryData(
			context.Background(),
			&backend.QueryDataRequest{
				Queries: []backend.DataQuery{
					{
						RefID: "A",
						JSON:  []byte(`{"zone":"entire-state"}`),
					},
				},
			},
		)

		frame := res.Responses["A"].Frames[0]

		require.Equal(t, frame.Fields[0].Name, "name")
		require.Equal(t, frame.At(0, 0).(string), "zone 1")
		require.Equal(t, frame.At(0, 1).(string), "zone 2")

		require.Equal(t, frame.Fields[1].Name, "rating")
		require.Equal(t, frame.At(1, 0).(int64), int64(1))
		require.Equal(t, frame.At(1, 1).(int64), int64(3))
	})

	t.Run("returns the specified zone", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.zones = []caic.Zone{
			{ID: "zone-1", Name: "zone 1", Rating: 1},
			{ID: "zone-2", Name: "zone 2", Rating: 3},
			{ID: "zone-3", Name: "zone 3", Rating: 4},
		}

		opts := plugin.DatasourceOpts(im)
		res, _ := opts.QueryDataHandler.QueryData(
			context.Background(),
			&backend.QueryDataRequest{
				Queries: []backend.DataQuery{
					{
						RefID: "A",
						JSON:  []byte(`{"zone":"zone-2"}`),
					},
				},
			},
		)

		frame := res.Responses["A"].Frames[0]
		require.Equal(t, 1, frame.Fields[0].Len())

		require.Equal(t, "name", frame.Fields[0].Name)
		require.Equal(t, "zone 2", frame.At(0, 0).(string))

		require.Equal(t, "rating", frame.Fields[1].Name)
		require.Equal(t, int64(3), frame.At(1, 0).(int64))
	})

	t.Run("returns different zones for different queries", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.zones = []caic.Zone{
			{ID: "zone-1", Name: "zone 1", Rating: 1},
			{ID: "zone-2", Name: "zone 2", Rating: 3},
			{ID: "zone-3", Name: "zone 3", Rating: 4},
		}

		opts := plugin.DatasourceOpts(im)
		res, _ := opts.QueryDataHandler.QueryData(
			context.Background(),
			&backend.QueryDataRequest{
				Queries: []backend.DataQuery{
					{
						RefID: "A",
						JSON:  []byte(`{"zone":"zone-2"}`),
					},
					{
						RefID: "B",
						JSON:  []byte(`{"zone":"zone-3"}`),
					},
				},
			},
		)

		frame := res.Responses["A"].Frames[0]
		require.Equal(t, 1, frame.Fields[0].Len())
		require.Equal(t, "zone 2", frame.At(0, 0).(string))

		frame = res.Responses["B"].Frames[0]
		require.Equal(t, 1, frame.Fields[0].Len())
		require.Equal(t, "zone 3", frame.At(0, 0).(string))
	})

	t.Run("return an error if it can't get zones", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.err = errors.New("something bad")
		opts := plugin.DatasourceOpts(im)
		_, err := opts.QueryDataHandler.QueryData(
			context.Background(),
			&backend.QueryDataRequest{
				Queries: []backend.DataQuery{
					{RefID: "A"},
				},
			},
		)

		require.EqualError(t, err, "something bad")
	})

	t.Run("returns returns an error if the request has bad json", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.zones = []caic.Zone{
			{ID: "zone-1", Name: "zone 1", Rating: 1},
			{ID: "zone-2", Name: "zone 2", Rating: 3},
			{ID: "zone-3", Name: "zone 3", Rating: 4},
		}

		opts := plugin.DatasourceOpts(im)
		_, err := opts.QueryDataHandler.QueryData(
			context.Background(),
			&backend.QueryDataRequest{
				Queries: []backend.DataQuery{
					{
						RefID: "A",
						JSON:  []byte(`{"zone":2}`), //can't marshal into string
					},
				},
			},
		)

		require.EqualError(t, err, "json: cannot unmarshal number into Go struct field .zone of type string")
	})
}
func TestCheckHealthHandler(t *testing.T) {
	t.Run("HealthStatusOK when can connect", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.canConnect = true

		opts := plugin.DatasourceOpts(im)
		res, _ := opts.CheckHealthHandler.CheckHealth(
			context.Background(),
			&backend.CheckHealthRequest{},
		)

		require.Equal(t, res.Status, backend.HealthStatusOk)
		require.Equal(t, res.Message, "Data source is working")
	})

	t.Run("HealthStatusError when can't connect", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.canConnect = false

		opts := plugin.DatasourceOpts(im)
		res, _ := opts.CheckHealthHandler.CheckHealth(
			context.Background(),
			&backend.CheckHealthRequest{},
		)

		require.Equal(t, res.Status, backend.HealthStatusError)
		require.Equal(t, res.Message, "Error reaching CAIC site")
	})

	t.Run("HealthStatusError instancemanager fails", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.err = errors.New("something bad")

		opts := plugin.DatasourceOpts(im)
		res, _ := opts.CheckHealthHandler.CheckHealth(
			context.Background(),
			&backend.CheckHealthRequest{},
		)

		require.Equal(t, res.Status, backend.HealthStatusError)
		require.Equal(t, res.Message, "something bad")
	})
}

func newFakeInstanceManager() *fakeInstanceManager {
	return &fakeInstanceManager{
		client: &fakeCaicClient{},
	}
}

type fakeInstanceManager struct {
	client *fakeCaicClient
	err    error
}

func (im *fakeInstanceManager) Get(pc backend.PluginContext) (instancemgmt.Instance, error) {
	return &plugin.CaicDatasource{
		Client: im.client,
	}, im.err
}

func (im *fakeInstanceManager) Do(pc backend.PluginContext, fn instancemgmt.InstanceCallbackFunc) error {
	return nil
}

type fakeCaicClient struct {
	canConnect bool
	zones      []caic.Zone
	err        error
}

func (c *fakeCaicClient) CanConnect() bool {
	return c.canConnect
}

func (c *fakeCaicClient) StateSummary() ([]caic.Zone, error) {
	return c.zones, c.err
}
