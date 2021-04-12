package caic_test

import (
	"net/http"
	"testing"

	"github.com/grafana/caic-datasource/pkg/caic"
	"github.com/stretchr/testify/require"
)

func TestGetRegionAspectInfo(t *testing.T) {
	t.Run("it returns whether or not each aspect is a danger by elevations", func(t *testing.T) {
		tc := setup(http.StatusOK, nil)
		tc.fakeHttp.resp <- avalancheProblem

		aspectDanger, _ := tc.caicClient.RegionAspectDanger(caic.SteamboatFlatTops)
		require.Equal(t, baseURL+"/caic/pub_bc_avo.php?zone_id=0", tc.fakeHttp.reqs[0].URL.String())
		require.Equal(t, http.MethodGet, tc.fakeHttp.reqs[0].Method)

		require.Equal(
			t,
			caic.AspectDanger{
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
			},
			aspectDanger,
		)
	})

	t.Run("it returns an error when the request fails", func(t *testing.T) {
		tc := setup(http.StatusNotFound, nil)

		_, err := tc.caicClient.RegionAspectDanger(caic.SteamboatFlatTops)
		require.NotNil(t, err)
	})
}

var (
	avalancheProblem = `
	<div class="ProblemRose">
		<div id="NBtl_0"  class="NBtl  off"></div>
		<div id="NTln_0"  class="NTln  on"></div>
		<div id="NAlp_0"  class="NAlp  on"></div>
		<div id="SAlp_0"  class="SAlp  off"></div>
		<div id="STln_0"  class="STln  off"></div>
		<div id="SBtl_0"  class="SBtl  off"></div>
		
		<div id="WBtl_0"  class="WBtl  off"></div>
		<div id="WTln_0"  class="WTln  off"></div>
		<div id="WAlp_0"  class="WAlp  off"></div>
		<div id="EAlp_0"  class="EAlp  off"></div>
		<div id="ETln_0"  class="ETln  off"></div>
		<div id="EBtl_0"  class="EBtl  off"></div>
		
		<div id="NWBtl_0" class="NWBtl off"></div>
		<div id="NWTln_0" class="NWTln on"></div>
		<div id="NWAlp_0" class="NWAlp on"></div>
		<div id="SEAlp_0" class="SEAlp off"></div>
		<div id="SETln_0" class="SETln off"></div>
		<div id="SEBtl_0" class="SEBtl off"></div>
		
		<div id="NEBtl_0" class="NEBtl off"></div>
		<div id="NETln_0" class="NETln on"></div>
		<div id="NEAlp_0" class="NEAlp on"></div>
		<div id="SWAlp_0" class="SWAlp off"></div>
		<div id="SWTln_0" class="SWTln off"></div>
		<div id="SWBtl_0" class="SWBtl off"></div>
	</div>`
)
