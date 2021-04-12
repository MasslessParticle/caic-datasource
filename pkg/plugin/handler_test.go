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

func TestQueryForZones(t *testing.T) {
	t.Run("returns all zones when zone is caic.EntireState", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.zones = []caic.Zone{
			{Index: 1, Name: "zone 1", Rating: 1},
			{Index: 2, Name: "zone 2", Rating: 3},
		}

		opts := plugin.DatasourceOpts(im)
		res, _ := opts.QueryDataHandler.QueryData(
			context.Background(),
			&backend.QueryDataRequest{
				Queries: []backend.DataQuery{
					{
						RefID: "A",
						JSON:  []byte(`{"zone":-1}`),
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

	t.Run("returns the specified zone with aspect", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.singleZone <- caic.Zone{
			Index:         2,
			Name:          "Zone 2",
			Rating:        4,
			AboveTreeline: 2,
			NearTreeline:  1,
			BelowTreeline: 4,
		}

		opts := plugin.DatasourceOpts(im)
		res, _ := opts.QueryDataHandler.QueryData(
			context.Background(),
			&backend.QueryDataRequest{
				Queries: []backend.DataQuery{
					{
						RefID: "A",
						JSON:  []byte(`{"zone":2}`),
					},
				},
			},
		)

		frame := res.Responses["A"].Frames[0]
		require.Equal(t, 1, frame.Fields[0].Len())

		require.Equal(t, "name", frame.Fields[0].Name)
		require.Equal(t, "Zone 2", frame.At(0, 0).(string))

		require.Equal(t, "rating", frame.Fields[1].Name)
		require.Equal(t, int64(4), frame.At(1, 0).(int64))

		require.Equal(t, "aboveTreeline", frame.Fields[2].Name)
		require.Equal(t, int64(2), frame.At(2, 0).(int64))

		require.Equal(t, "nearTreeline", frame.Fields[3].Name)
		require.Equal(t, int64(1), frame.At(3, 0).(int64))

		require.Equal(t, "belowTreeline", frame.Fields[4].Name)
		require.Equal(t, int64(4), frame.At(4, 0).(int64))
	})

	t.Run("returns different zones for different queries", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.singleZone <- caic.Zone{Index: 2, Name: "Zone 2", Rating: 3}
		im.client.singleZone <- caic.Zone{Index: 3, Name: "Zone 3", Rating: 3}

		opts := plugin.DatasourceOpts(im)
		res, _ := opts.QueryDataHandler.QueryData(
			context.Background(),
			&backend.QueryDataRequest{
				Queries: []backend.DataQuery{
					{
						RefID: "A",
						JSON:  []byte(`{"zone":2}`),
					},
					{
						RefID: "B",
						JSON:  []byte(`{"zone":3}`),
					},
				},
			},
		)

		frame := res.Responses["A"].Frames[0]
		require.Equal(t, 1, frame.Fields[0].Len())
		require.Equal(t, "Zone 2", frame.At(0, 0).(string))

		frame = res.Responses["B"].Frames[0]
		require.Equal(t, 1, frame.Fields[0].Len())
		require.Equal(t, "Zone 3", frame.At(0, 0).(string))
	})

	t.Run("return an error if it can't get zones", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.singleZone <- caic.Zone{Index: 2, Name: "Zone 2", Rating: 3}
		im.client.err = errors.New("something bad")

		opts := plugin.DatasourceOpts(im)
		_, err := opts.QueryDataHandler.QueryData(
			context.Background(),
			&backend.QueryDataRequest{
				Queries: []backend.DataQuery{
					{
						RefID: "A",
						JSON:  []byte(`{"zone":2}`),
					},
				},
			},
		)

		require.Contains(t, err.Error(), "something bad")
	})

	t.Run("returns returns an error if the request has bad json", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.zones = []caic.Zone{
			{Index: 1, Name: "zone 1", Rating: 1},
			{Index: 2, Name: "zone 2", Rating: 3},
			{Index: 3, Name: "zone 3", Rating: 4},
		}

		opts := plugin.DatasourceOpts(im)
		_, err := opts.QueryDataHandler.QueryData(
			context.Background(),
			&backend.QueryDataRequest{
				Queries: []backend.DataQuery{
					{
						RefID: "A",
						JSON:  []byte(`{"zone": "2"}`), //can't marshal into int
					},
				},
			},
		)

		require.Contains(t, err.Error(), "json: cannot unmarshal string into Go struct field .zone of type caic.Region")
	})
}

func TestQueryForProblems(t *testing.T) {
	t.Run("it returns aspect problem data", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.singleZone <- caic.Zone{}
		im.client.aspectDanger = caic.AspectDanger{
			Region:        caic.SteamboatFlatTops,
			BelowTreeline: caic.OrdinalDanger{},
			NearTreeline: caic.OrdinalDanger{
				North:     true,
				NorthEast: true,
				NorthWest: true,
			},
			AboveTreeline: caic.OrdinalDanger{
				North:     true,
				NorthEast: true,
				NorthWest: true,
			},
		}

		opts := plugin.DatasourceOpts(im)
		res, _ := opts.QueryDataHandler.QueryData(
			context.Background(),
			&backend.QueryDataRequest{
				Queries: []backend.DataQuery{
					{
						RefID: "A",
						JSON:  []byte(`{"zone":2}`),
					},
				},
			},
		)

		frame := res.Responses["A"].Frames[1]

		require.Equal(t, "ordinals", frame.Fields[0].Name)
		expectedDirections := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
		for i := 0; i < frame.Fields[0].Len(); i++ {
			require.Equal(t, expectedDirections[i], frame.Fields[0].At(i).(string))
		}

		require.Equal(t, "degrees", frame.Fields[1].Name)
		expected := []int32{0, 45, 90, 135, 180, 225, 270, 315}
		for i := 0; i < frame.Fields[0].Len(); i++ {
			require.Equal(t, expected[i], frame.Fields[1].At(i).(int32))
		}

		require.Equal(t, "aboveTreeline", frame.Fields[2].Name)
		expected = []int32{1, 1, 0, 0, 0, 0, 0, 1}
		for i := 0; i < frame.Fields[0].Len(); i++ {
			require.Equal(t, expected[i], frame.Fields[2].At(i).(int32))
		}

		require.Equal(t, "nearTreeline", frame.Fields[3].Name)
		expected = []int32{1, 1, 0, 0, 0, 0, 0, 1}
		for i := 0; i < frame.Fields[0].Len(); i++ {
			require.Equal(t, expected[i], frame.Fields[3].At(i).(int32))
		}

		require.Equal(t, "belowTreeline", frame.Fields[4].Name)
		expected = []int32{0, 0, 0, 0, 0, 0, 0, 0}
		for i := 0; i < frame.Fields[0].Len(); i++ {
			require.Equal(t, expected[i], frame.Fields[4].At(i).(int32))
		}
	})

	t.Run("it doesn't return anything if region is EntireState", func(t *testing.T) {
		im := newFakeInstanceManager()
		im.client.singleZone <- caic.Zone{}

		opts := plugin.DatasourceOpts(im)
		res, _ := opts.QueryDataHandler.QueryData(
			context.Background(),
			&backend.QueryDataRequest{
				Queries: []backend.DataQuery{
					{
						RefID: "A",
						JSON:  []byte(`{"zone":-1}`),
					},
				},
			},
		)

		require.Len(t, res.Responses["A"].Frames, 1)
		require.Equal(t, "Zones", res.Responses["A"].Frames[0].Name)
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
		client: &fakeCaicClient{
			singleZone:   make(chan caic.Zone, 10),
			aspectDanger: caic.AspectDanger{},
		},
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
	canConnect   bool
	zones        []caic.Zone
	aspectDanger caic.AspectDanger
	singleZone   chan caic.Zone
	err          error
}

func (c *fakeCaicClient) CanConnect() bool {
	return c.canConnect
}

func (c *fakeCaicClient) StateSummary() ([]caic.Zone, error) {
	return c.zones, c.err
}

func (c *fakeCaicClient) RegionSummary(r caic.Region) (caic.Zone, error) {
	select {
	case zone := <-c.singleZone:
		return zone, c.err
	default:
		panic("called without any responses setup")
	}
}

func (c *fakeCaicClient) RegionAspectDanger(caic.Region) (caic.AspectDanger, error) {
	return c.aspectDanger, c.err
}
