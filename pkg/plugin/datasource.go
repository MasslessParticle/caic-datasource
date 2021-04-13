package plugin

import "github.com/grafana/caic-datasource/pkg/caic"

type caicClient interface {
	CanConnect() bool
	Summary(caic.Region) ([]caic.Zone, error)
	AspectDanger(caic.Region) (caic.AspectDanger, error)
}

type CaicDatasource struct {
	Client caicClient
}

func (s *CaicDatasource) Dispose() {
	// Called before creatinga a new instance to allow plugin authors
	// to cleanup.
}
