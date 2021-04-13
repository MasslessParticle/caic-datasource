package cache

import (
	"sync"
	"time"

	"github.com/grafana/caic-datasource/pkg/caic"
)

type client interface {
	CanConnect() bool
	Summary(caic.Region) ([]caic.Zone, error)
	AspectDanger(caic.Region) (caic.AspectDanger, error)
}

type zone struct {
	t time.Time
	z []caic.Zone
}

type aspectDanger struct {
	t  time.Time
	ad caic.AspectDanger
}

type Cache struct {
	m                 sync.Mutex
	client            client
	regionCache       map[string]zone
	aspectDangerCache map[string]aspectDanger
	cacheDuration     time.Duration
}

func NewCaicClientCache(c client, opts ...CacheOption) *Cache {
	cache := &Cache{
		client:            c,
		regionCache:       make(map[string]zone),
		aspectDangerCache: make(map[string]aspectDanger),
		cacheDuration:     time.Hour,
	}

	for _, o := range opts {
		o(cache)
	}

	return cache
}

type CacheOption func(c *Cache)

func WithCacheDuration(d time.Duration) CacheOption {
	return func(c *Cache) {
		c.cacheDuration = d
	}

}

func (c *Cache) Summary(r caic.Region) ([]caic.Zone, error) {
	c.m.Lock()
	defer c.m.Unlock()

	cached, ok := c.regionCache[r.String()]
	if ok && time.Since(cached.t) < c.cacheDuration {
		return cached.z, nil
	}

	z, err := c.client.Summary(r)
	if err != nil {
		return nil, err
	}

	c.regionCache[r.String()] = zone{
		t: time.Now(),
		z: z,
	}

	return z, nil
}

func (c *Cache) AspectDanger(r caic.Region) (caic.AspectDanger, error) {
	c.m.Lock()
	defer c.m.Unlock()

	ad, ok := c.aspectDangerCache[r.String()]
	if ok && time.Since(ad.t) < c.cacheDuration {
		return ad.ad, nil
	}

	a, err := c.client.AspectDanger(r)
	if err != nil {
		return caic.AspectDanger{}, err
	}

	c.aspectDangerCache[r.String()] = aspectDanger{
		t:  time.Now(),
		ad: a,
	}

	return a, nil
}

func (c *Cache) CanConnect() bool {
	return c.client.CanConnect()
}
