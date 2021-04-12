package plugin

import "github.com/grafana/caic-datasource/pkg/caic"

type CaicClient interface {
	CanConnect() bool
	StateSummary() ([]caic.Zone, error)
	RegionSummary(caic.Region) (caic.Zone, error)
	RegionAspectDanger(caic.Region) (caic.AspectDanger, error)
}

type CaicDatasource struct {
	Client CaicClient
}

func (s *CaicDatasource) Dispose() {
	// Called before creatinga a new instance to allow plugin authors
	// to cleanup.
}
