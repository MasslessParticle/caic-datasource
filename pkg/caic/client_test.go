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

func TestClientGetsStateSummary(t *testing.T) {
	t.Run("it returns an array of state zones", func(t *testing.T) {
		tc := setup(http.StatusOK, nil)
		tc.fakeHttp.resp <- homePage

		expected := []caic.Zone{
			{ID: "zone-0", Name: "Zone 0", Rating: 3},
			{ID: "zone-1", Name: "Zone 1", Rating: 2},
			{ID: "zone-12", Name: "Zone 12", Rating: 1},
		}

		zones, _ := tc.caicClient.StateSummary()

		require.Equal(t, baseURL+"/caic/fx_map.php", tc.fakeHttp.reqs[0].URL.String())
		require.Equal(t, http.MethodGet, tc.fakeHttp.reqs[0].Method)
		require.Equal(t, expected, zones)
	})

	t.Run("it returns an error if the CAIC website can't be reached", func(t *testing.T) {
		tc := setup(http.StatusNotFound, nil)

		_, err := tc.caicClient.StateSummary()
		require.NotNil(t, err)
	})
}

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
		tc.fakeHttp.resp <- forecastContainer
		tc.fakeHttp.resp <- forecastFragment

		zone, _ := tc.caicClient.RegionSummary("front-range")

		require.Equal(t, baseURL+"/forecasts/backcountry-avalanche/front-range/", tc.fakeHttp.reqs[0].URL.String())
		require.Equal(t, http.MethodGet, tc.fakeHttp.reqs[0].Method)

		require.Equal(t, baseURL+"/caic/pub_bc_avo.php?zone_id=0", tc.fakeHttp.reqs[1].URL.String())
		require.Equal(t, http.MethodGet, tc.fakeHttp.reqs[1].Method)

		require.Equal(
			t,
			zone,
			caic.Zone{
				ID:            "front-range",
				Name:          "Front Range",
				Rating:        4,
				AboveTreeline: 3,
				NearTreeline:  2,
				BelowTreeline: 4,
			})
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
	baseURL = "http://www.caic-url.com"

	homePage = `
zone[0]='Zone 0';
url[0]='http://caic-url.com/forecasts/backcountry-avalanche/zone-0/';
rating[0]=3;
--
zone[1]='Zone 1';
url[1]='http://caic-url.com/forecasts/backcountry-avalanche/zone-1/';
rating[1]=2;
--
zone[12]='Zone 12';
url[12]='http://caic-url.com/forecasts/backcountry-avalanche/zone-12/';
rating[12]=1;
`

	forecastContainer = `
<head>
	<title>Front Range</title>
</head>
<body>
	<div class="site-container">
		<div class="site-inner">
			<div>
				<div>
					<main>
						<article>
							<div>
								<iframe src="/caic/pub_bc_avo.php?zone_id=0"></iframe>
							</div>
						</article>
					</main>
				</div>	
			</div>
		</div>
	</div>
</body>
`

	forecastFragment = `
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
