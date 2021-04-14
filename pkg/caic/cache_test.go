package caic_test

import (
	"errors"
	"testing"
	"time"

	"github.com/grafana/caic-datasource/pkg/caic"
	"github.com/stretchr/testify/require"
)

func TestSummary(t *testing.T) {
	t.Run("it caches responses for duration", func(t *testing.T) {
		client := newFakeClient()
		client.regionResponse <- []caic.Zone{{Name: "Zone 1"}}
		client.regionResponse <- []caic.Zone{{Name: "Zone 2"}}

		cache := caic.NewClientCache(client, caic.WithCacheDuration(10*time.Millisecond))

		call, err := cache.Summary(caic.SteamboatFlatTops)
		require.Nil(t, err)

		cachedCall, err := cache.Summary(caic.SteamboatFlatTops)
		require.Nil(t, err)

		time.Sleep(20 * time.Millisecond)

		secondCall, err := cache.Summary(caic.SteamboatFlatTops)
		require.Nil(t, err)

		require.Equal(t, call, cachedCall)
		require.Equal(t, "Zone 1", call[0].Name)
		require.Equal(t, "Zone 2", secondCall[0].Name)
	})

	// If incorrect, this test will fail when run with go test -race
	t.Run("it is threadsafe", func(t *testing.T) {
		client := newFakeClient()
		cache := caic.NewClientCache(client, caic.WithCacheDuration(10*time.Millisecond))

		start := make(chan struct{})
		stop := make(chan struct{})

		go readRegions(start, stop, cache)
		go readRegions(start, stop, cache)

		close(start)
		time.Sleep(10 * time.Millisecond)
		close(stop)
	})

	t.Run("it doesn't cache errors", func(t *testing.T) {
		client := newFakeClient()
		client.regionResponse <- nil
		client.regionResponse <- []caic.Zone{{Name: "Zone 2"}}
		client.err <- errors.New("something bad")

		cache := caic.NewClientCache(client, caic.WithCacheDuration(10*time.Millisecond))

		_, err := cache.Summary(caic.SteamboatFlatTops)
		require.NotNil(t, err)

		secondCall, err := cache.Summary(caic.SteamboatFlatTops)
		require.Nil(t, err)

		require.Equal(t, "Zone 2", secondCall[0].Name)
	})
}

func TestAspectDangerSummary(t *testing.T) {
	t.Run("it caches responses for duration", func(t *testing.T) {
		client := newFakeClient()
		client.aspectDangerResponse <- caic.AspectDanger{Region: caic.SteamboatFlatTops}
		client.aspectDangerResponse <- caic.AspectDanger{Region: caic.SawatchRange}

		cache := caic.NewClientCache(client, caic.WithCacheDuration(10*time.Millisecond))

		call, err := cache.AspectDanger(caic.SteamboatFlatTops)
		require.Nil(t, err)

		cachedCall, err := cache.AspectDanger(caic.SteamboatFlatTops)
		require.Nil(t, err)

		time.Sleep(20 * time.Millisecond)

		secondCall, err := cache.AspectDanger(caic.SteamboatFlatTops)
		require.Nil(t, err)

		require.Equal(t, call, cachedCall)
		require.Equal(t, caic.SteamboatFlatTops, call.Region)
		require.Equal(t, caic.SawatchRange, secondCall.Region)
	})

	//If incorrect, this test will fail when run with go test -race
	t.Run("it is threadsafe", func(t *testing.T) {
		client := newFakeClient()
		cache := caic.NewClientCache(client, caic.WithCacheDuration(10*time.Millisecond))

		start := make(chan struct{})
		stop := make(chan struct{})

		go readAspectDanger(start, stop, cache)
		go readAspectDanger(start, stop, cache)

		close(start)
		time.Sleep(10 * time.Millisecond)
		close(stop)
	})

	t.Run("it doesn't cache errors", func(t *testing.T) {
		client := newFakeClient()
		client.aspectDangerResponse <- caic.AspectDanger{}
		client.aspectDangerResponse <- caic.AspectDanger{Region: caic.Aspen}
		client.err <- errors.New("something bad")

		cache := caic.NewClientCache(client, caic.WithCacheDuration(10*time.Millisecond))

		_, err := cache.AspectDanger(caic.SteamboatFlatTops)
		require.NotNil(t, err)

		secondCall, err := cache.AspectDanger(caic.SteamboatFlatTops)
		require.Nil(t, err)

		require.Equal(t, caic.Aspen, secondCall.Region)
	})
}

func TestCanConnect(t *testing.T) {
	t.Run("it does not cache responses", func(t *testing.T) {
		client := newFakeClient()
		client.canConnectResponse <- true
		client.canConnectResponse <- false

		cache := caic.NewClientCache(client)
		require.True(t, cache.CanConnect())
		require.False(t, cache.CanConnect())
	})
}

func readRegions(start, stop chan struct{}, c *caic.Cache) {
	<-start
	for {
		select {
		case <-stop:
			return
		default:
			c.Summary(caic.SteamboatFlatTops)
		}
	}
}

func readAspectDanger(start, stop chan struct{}, c *caic.Cache) {
	<-start
	for {
		select {
		case <-stop:
			return
		default:
			c.AspectDanger(caic.SteamboatFlatTops)
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

func (c *fakeClient) Summary(caic.Region) ([]caic.Zone, error) {
	select {
	case ret := <-c.regionResponse:
		return ret, c.error()
	default:
		return nil, c.error()
	}
}

func (c *fakeClient) AspectDanger(caic.Region) (caic.AspectDanger, error) {
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
