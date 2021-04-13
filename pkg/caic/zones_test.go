package caic_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/grafana/caic-datasource/pkg/caic"
	"github.com/stretchr/testify/require"
)

func TestClientCanConnect(t *testing.T) {
	t.Run("it returns true when it can connect", func(t *testing.T) {
		tc := setup(http.StatusOK, nil)
		require.True(t, tc.caicClient.CanConnect())
	})

	t.Run("return false when it gets a non 200", func(t *testing.T) {
		tc := setup(http.StatusBadGateway, nil)
		require.False(t, tc.caicClient.CanConnect())
	})

	t.Run("return false when the client has an error", func(t *testing.T) {
		tc := setup(http.StatusOK, errors.New("something bad happened"))
		require.False(t, tc.caicClient.CanConnect())
	})
}

func TestGetRegionSummary(t *testing.T) {
	t.Run("returns the forecast by elevation for a single zone", func(t *testing.T) {
		tc := setup(http.StatusOK, nil)
		tc.fakeHttp.resp <- forecast

		zone, _ := tc.caicClient.RegionSummary(caic.SteamboatFlatTops)
		require.Equal(t, baseURL+"/caic/pub_bc_avo.php?zone_id=0", tc.fakeHttp.reqs[0].URL.String())
		require.Equal(t, http.MethodGet, tc.fakeHttp.reqs[0].Method)

		require.Equal(
			t,
			zone,
			[]caic.Zone{
				{
					Index:         0,
					Name:          caic.SteamboatFlatTops.String(),
					Rating:        4,
					AboveTreeline: 3,
					NearTreeline:  2,
					BelowTreeline: 4,
				},
			})
	})

	t.Run("it returns an array of state zones when region is EntireState", func(t *testing.T) {
		tc := setup(http.StatusOK, nil)
		tc.fakeHttp.resp <- forecast
		tc.fakeHttp.resp <- forecast
		tc.fakeHttp.resp <- forecast
		tc.fakeHttp.resp <- forecast
		tc.fakeHttp.resp <- forecast
		tc.fakeHttp.resp <- forecast
		tc.fakeHttp.resp <- forecast
		tc.fakeHttp.resp <- forecast
		tc.fakeHttp.resp <- forecast
		tc.fakeHttp.resp <- forecast

		expected := []caic.Zone{
			{Index: caic.SteamboatFlatTops, Name: caic.SteamboatFlatTops.String(), Rating: 4, AboveTreeline: 3, NearTreeline: 2, BelowTreeline: 4},
			{Index: caic.FrontRange, Name: caic.FrontRange.String(), Rating: 4, AboveTreeline: 3, NearTreeline: 2, BelowTreeline: 4},
			{Index: caic.VailSummitCounty, Name: caic.VailSummitCounty.String(), Rating: 4, AboveTreeline: 3, NearTreeline: 2, BelowTreeline: 4},
			{Index: caic.SawatchRange, Name: caic.SawatchRange.String(), Rating: 4, AboveTreeline: 3, NearTreeline: 2, BelowTreeline: 4},
			{Index: caic.Aspen, Name: caic.Aspen.String(), Rating: 4, AboveTreeline: 3, NearTreeline: 2, BelowTreeline: 4},
			{Index: caic.Gunnison, Name: caic.Gunnison.String(), Rating: 4, AboveTreeline: 3, NearTreeline: 2, BelowTreeline: 4},
			{Index: caic.GrandMesa, Name: caic.GrandMesa.String(), Rating: 4, AboveTreeline: 3, NearTreeline: 2, BelowTreeline: 4},
			{Index: caic.NorthernSanJuan, Name: caic.NorthernSanJuan.String(), Rating: 4, AboveTreeline: 3, NearTreeline: 2, BelowTreeline: 4},
			{Index: caic.SouthernSanJuan, Name: caic.SouthernSanJuan.String(), Rating: 4, AboveTreeline: 3, NearTreeline: 2, BelowTreeline: 4},
			{Index: caic.SangreDeCristo, Name: caic.SangreDeCristo.String(), Rating: 4, AboveTreeline: 3, NearTreeline: 2, BelowTreeline: 4},
		}

		zones, _ := tc.caicClient.RegionSummary(caic.EntireState)

		require.Equal(t, baseURL+"/caic/pub_bc_avo.php?zone_id=0", tc.fakeHttp.reqs[0].URL.String())
		require.Equal(t, baseURL+"/caic/pub_bc_avo.php?zone_id=9", tc.fakeHttp.reqs[9].URL.String())
		require.Equal(t, http.MethodGet, tc.fakeHttp.reqs[0].Method)

		require.Equal(t, expected, zones)
	})

	t.Run("it returns an error if the CAIC website can't be reached", func(t *testing.T) {
		tc := setup(http.StatusNotFound, nil)

		_, err := tc.caicClient.RegionSummary(caic.EntireState)
		require.NotNil(t, err)
	})
}

type testContext struct {
	fakeHttp   *spyHttpClient
	caicClient *caic.Client
}

func setup(responseCode int, httpError error) testContext {
	fakeHttp := &spyHttpClient{
		resp:     make(chan string, 10),
		respCode: responseCode,
		err:      httpError,
	}

	return testContext{
		fakeHttp:   fakeHttp,
		caicClient: caic.NewClient(baseURL, fakeHttp),
	}
}

type spyHttpClient struct {
	reqs     []*http.Request
	resp     chan string
	respCode int
	err      error
}

func (c *spyHttpClient) Do(r *http.Request) (*http.Response, error) {
	c.reqs = append(c.reqs, r)

	var resp string
	select {
	case resp = <-c.resp:
	default:
	}

	return &http.Response{
		StatusCode: c.respCode,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(resp))),
	}, c.err
}

var (
	baseURL  = "http://www.caic-url.com"
	forecast = `
<div id="avalanche-forecast">
	<table class="table table-striped-body table-treeline">
		<tbody>
			<tr>
				<td class="today-text above_danger_low" style="">
						<strong>Considerable (3)</strong>
				</td>
				<td class="today-text tomorrow_danger_low">
						<strong>Low (1)</strong>
				</td>
			</tr>
			<tr>
				<td class="today-text near_danger_low">
						<strong>Moderate (2)</strong>
				</td>
				<td class="today-text tomorrow_danger_low">
						<strong>Low (1)</strong>
				</td>
			</tr>
			<tr>
				<td class="today-text below_danger_moderate">
						<strong>High (4)</strong>
				</td>
				<td class="today-text tomorrow_danger_low">
						<strong>Low (1)</strong>
				</td>
			</tr>
		</tbody>
	</table>
</div>
`
)
