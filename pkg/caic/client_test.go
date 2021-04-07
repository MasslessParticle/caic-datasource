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
	fakeHttp := newSpyHttpClient()
	client := caic.NewClient(baseURL, fakeHttp)

	t.Run("it calls the right url for the whole state", func(t *testing.T) {
		client.StateSummary()
		require.Equal(t, baseURL+"/caic/fx_map.php", fakeHttp.req.URL.String())
		require.Equal(t, http.MethodGet, fakeHttp.req.Method)
	})

	t.Run("it returns an array of state zones", func(t *testing.T) {
		fakeHttp.resp = mockCaicPage
		zones, _ := client.StateSummary()

		expected := []caic.Zone{
			{"Zone 0", "http://caic-url.com/zone_0", "3"},
			{"Zone 1", "http://caic-url.com/zone_1", "2"},
			{"Zone 12", "http://caic-url.com/zone_12", "1"},
		}

		require.Equal(t, zones, expected)
	})

	t.Run("it returns an array of state zones", func(t *testing.T) {
		fakeHttp.resp = mockCaicPage
		zones, _ := client.StateSummary()

		expected := []caic.Zone{
			{"Zone 0", "http://caic-url.com/zone_0", "3"},
			{"Zone 1", "http://caic-url.com/zone_1", "2"},
			{"Zone 12", "http://caic-url.com/zone_12", "1"},
		}

		require.Equal(t, zones, expected)
	})

	t.Run("it returns an error if the CAIC website can't be reached", func(t *testing.T) {
		fakeHttp.respCode = http.StatusNotFound
		_, err := client.StateSummary()

		require.NotNil(t, err)
	})
}

func TestClientCanConnect(t *testing.T) {
	fakeHttp := newSpyHttpClient()
	client := caic.NewClient(baseURL, fakeHttp)

	t.Run("it returns true when it can connect", func(t *testing.T) {
		require.True(t, client.CanConnect())
	})

	t.Run("return false when it gets a non 200", func(t *testing.T) {
		fakeHttp.respCode = http.StatusBadGateway
		require.False(t, client.CanConnect())
	})

	t.Run("return false when the client has an error", func(t *testing.T) {
		fakeHttp.respCode = http.StatusOK
		fakeHttp.err = errors.New("something bad happened")
		require.False(t, client.CanConnect())
	})
}

type spyHttpClient struct {
	req      *http.Request
	resp     string
	respCode int
	err      error
}

func newSpyHttpClient() *spyHttpClient {
	return &spyHttpClient{
		respCode: 200,
	}
}

func (c *spyHttpClient) Do(r *http.Request) (*http.Response, error) {
	c.req = r
	return &http.Response{
		StatusCode: c.respCode,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(c.resp))),
	}, c.err
}

var baseURL = "http://www.caic-url.com"
var mockCaicPage = `
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
