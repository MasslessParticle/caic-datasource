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
		tc := setup(homePage, http.StatusOK, nil)

		expected := []caic.Zone{
			{0, "Zone 0", "http://caic-url.com/zone_0", 3},
			{1, "Zone 1", "http://caic-url.com/zone_1", 2},
			{12, "Zone 12", "http://caic-url.com/zone_12", 1},
		}

		zones, _ := tc.caicClend.StateSummary()

		require.Equal(t, baseURL+"/caic/fx_map.php", tc.fakeHttp.req.URL.String())
		require.Equal(t, http.MethodGet, tc.fakeHttp.req.Method)
		require.Equal(t, zones, expected)
	})

	t.Run("it returns an error if the CAIC website can't be reached", func(t *testing.T) {
		tc := setup("", http.StatusNotFound, nil)

		_, err := tc.caicClend.StateSummary()
		require.NotNil(t, err)
	})
}

func TestClientCanConnect(t *testing.T) {
	t.Run("it returns true when it can connect", func(t *testing.T) {
		tc := setup("", http.StatusOK, nil)
		require.True(t, tc.caicClend.CanConnect())
	})

	t.Run("return false when it gets a non 200", func(t *testing.T) {
		tc := setup("", http.StatusBadGateway, nil)
		require.False(t, tc.caicClend.CanConnect())
	})

	t.Run("return false when the client has an error", func(t *testing.T) {
		tc := setup("", http.StatusOK, errors.New("something bad happened"))
		require.False(t, tc.caicClend.CanConnect())
	})
}

type testContext struct {
	fakeHttp  *spyHttpClient
	caicClend *caic.Client
}

func setup(response string, responseCode int, httpError error) testContext {
	fakeHttp := &spyHttpClient{
		resp:     response,
		respCode: responseCode,
		err:      httpError,
	}

	return testContext{
		fakeHttp:  fakeHttp,
		caicClend: caic.NewClient(baseURL, fakeHttp),
	}
}

type spyHttpClient struct {
	req      *http.Request
	resp     string
	respCode int
	err      error
}

func (c *spyHttpClient) Do(r *http.Request) (*http.Response, error) {
	c.req = r
	return &http.Response{
		StatusCode: c.respCode,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(c.resp))),
	}, c.err
}

var (
	baseURL = "http://www.caic-url.com"

	homePage = `
zone[0]='Zone 0';
url[0]='http://caic-url.com/zone_0';
rating[0]=3;
--
zone[1]='Zone 1';
url[1]='http://caic-url.com/zone_1';
rating[1]=2;
--
zone[12]='Zone 12';
url[12]='http://caic-url.com/zone_12';
rating[12]=1;
`
)
