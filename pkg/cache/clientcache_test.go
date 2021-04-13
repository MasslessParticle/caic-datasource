package cache_test

import (
	"errors"
	"testing"
	"time"

	"github.com/grafana/caic-datasource/pkg/cache"
	"github.com/grafana/caic-datasource/pkg/caic"
	"github.com/stretchr/testify/require"
)

func TestRegionSummary(t *testing.T) {
	t.Run("it caches responses for duration", func(t *testing.T) {
		client := newFakeClient()
		client.regionResponse <- []caic.Zone{{Name: "Zone 1"}}
		client.regionResponse <- []caic.Zone{{Name: "Zone 2"}}

		c := cache.NewCaicClientCache(client, cache.WithCacheDuration(10*time.Millisecond))

		call, err := c.RegionSummary(caic.SteamboatFlatTops)
		require.Nil(t, err)

		cachedCall, err := c.RegionSummary(caic.SteamboatFlatTops)
		require.Nil(t, err)

		time.Sleep(20 * time.Millisecond)

		secondCall, err := c.RegionSummary(caic.SteamboatFlatTops)
		require.Nil(t, err)

		require.Equal(t, call, cachedCall)
		require.Equal(t, "Zone 1", call[0].Name)
		require.Equal(t, "Zone 2", secondCall[0].Name)
	})

	// If incorrect, this test will fail when run with go test -race
	t.Run("it is threadsafe", func(t *testing.T) {
		client := newFakeClient()
		c := cache.NewCaicClientCache(client, cache.WithCacheDuration(10*time.Millisecond))

		start := make(chan struct{})
		stop := make(chan struct{})

		go readRegions(start, stop, c)
		go readRegions(start, stop, c)

		close(start)
		time.Sleep(10 * time.Millisecond)
		close(stop)
	})

	t.Run("it doesn't cache errors", func(t *testing.T) {
		client := newFakeClient()
		client.regionResponse <- nil
		client.regionResponse <- []caic.Zone{{Name: "Zone 2"}}
		client.err <- errors.New("something bad")

		c := cache.NewCaicClientCache(client, cache.WithCacheDuration(10*time.Millisecond))

		_, err := c.RegionSummary(caic.SteamboatFlatTops)
		require.NotNil(t, err)

		secondCall, err := c.RegionSummary(caic.SteamboatFlatTops)
		require.Nil(t, err)

		require.Equal(t, "Zone 2", secondCall[0].Name)
	})
}

func TestAspectDangerSummary(t *testing.T) {
	t.Run("it caches responses for duration", func(t *testing.T) {
		client := newFakeClient()
		client.aspectDangerResponse <- caic.AspectDanger{Region: caic.SteamboatFlatTops}
		client.aspectDangerResponse <- caic.AspectDanger{Region: caic.SawatchRange}

		c := cache.NewCaicClientCache(client, cache.WithCacheDuration(10*time.Millisecond))

		call, err := c.RegionAspectDanger(caic.SteamboatFlatTops)
		require.Nil(t, err)

		cachedCall, err := c.RegionAspectDanger(caic.SteamboatFlatTops)
		require.Nil(t, err)

		time.Sleep(20 * time.Millisecond)

		secondCall, err := c.RegionAspectDanger(caic.SteamboatFlatTops)
		require.Nil(t, err)

		require.Equal(t, call, cachedCall)
		require.Equal(t, caic.SteamboatFlatTops, call.Region)
		require.Equal(t, caic.SawatchRange, secondCall.Region)
	})

	//If incorrect, this test will fail when run with go test -race
	t.Run("it is threadsafe", func(t *testing.T) {
		client := newFakeClient()
		c := cache.NewCaicClientCache(client, cache.WithCacheDuration(10*time.Millisecond))

		start := make(chan struct{})
		stop := make(chan struct{})

		go readAspectDanger(start, stop, c)
		go readAspectDanger(start, stop, c)

		close(start)
		time.Sleep(10 * time.Millisecond)
		close(stop)
	})

	t.Run("it doesn't cache errors", func(t *testing.T) {
		client := newFakeClient()
		client.aspectDangerResponse <- caic.AspectDanger{}
		client.aspectDangerResponse <- caic.AspectDanger{Region: caic.Aspen}
		client.err <- errors.New("something bad")

		c := cache.NewCaicClientCache(client, cache.WithCacheDuration(10*time.Millisecond))

		_, err := c.RegionAspectDanger(caic.SteamboatFlatTops)
		require.NotNil(t, err)

		secondCall, err := c.RegionAspectDanger(caic.SteamboatFlatTops)
		require.Nil(t, err)

		require.Equal(t, caic.Aspen, secondCall.Region)
	})
}

func TestCanConnect(t *testing.T) {
	t.Run("it does not cache responses", func(t *testing.T) {
		client := newFakeClient()
		client.canConnectResponse <- true
		client.canConnectResponse <- false

		c := cache.NewCaicClientCache(client)
		require.True(t, c.CanConnect())
		require.False(t, c.CanConnect())
	})
}

func readRegions(start, stop chan struct{}, c *cache.Cache) {
	<-start
	for {
		select {
		case <-stop:
			return
		default:
			c.RegionSummary(caic.SteamboatFlatTops)
		}
	}
}

func readAspectDanger(start, stop chan struct{}, c *cache.Cache) {
	<-start
	for {
		select {
		case <-stop:
			return
		default:
			c.RegionAspectDanger(caic.SteamboatFlatTops)
		}
	}
}

func newFakeClient() *fakeClient {
	return &fakeClient{
		regionResponse:       make(chan []caic.Zone, 10),
		aspectDangerResponse: make(chan caic.AspectDanger, 10),
		canConnectResponse:   make(chan bool, 10),
		err:                  make(chan error, 10),
	}
}

type fakeClient struct {
	regionResponse       chan []caic.Zone
	aspectDangerResponse chan caic.AspectDanger
	canConnectResponse   chan bool
	err                  chan error
}

func (c *fakeClient) CanConnect() bool {
	select {
	case ret := <-c.canConnectResponse:
		return ret
	default:
		return false
	}
}

func (c *fakeClient) RegionSummary(caic.Region) ([]caic.Zone, error) {
	select {
	case ret := <-c.regionResponse:
		return ret, c.error()
	default:
		return nil, c.error()
	}
}

func (c *fakeClient) RegionAspectDanger(caic.Region) (caic.AspectDanger, error) {
	select {
	case ret := <-c.aspectDangerResponse:
		return ret, c.error()
	default:
		return caic.AspectDanger{}, c.error()
	}
}

func (c *fakeClient) error() error {
	select {
	case err := <-c.err:
		return err
	default:
		return nil
	}
}
