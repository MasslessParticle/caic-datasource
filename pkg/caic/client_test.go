package caic_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/grafana/caic-datasource/pkg/caic"
	"github.com/stretchr/testify/require"
)

func TestClientGetsRegions(t *testing.T) {
	fakeHttp := newSpyHttpClient()
	client := caic.NewClient(baseURL, fakeHttp)

	t.Run("it calls the right url for the whole state", func(t *testing.T) {
		client.StateSummary()
		require.Equal(t, baseURL+"/caic/fx_map.php", fakeHttp.req.URL.String())
	})

	t.Run("it returns an array of state zones", func(t *testing.T) {
		fakeHttp.resp = mockCaicPage
		zones := client.StateSummary()

		expected := []caic.Zone{
			{"0", "Zone 0", "http://caic-url.com/zone_0", "3"},
			{"1", "Zone 1", "http://caic-url.com/zone_1", "2"},
			{"12", "Zone 12", "http://caic-url.com/zone_12", "1"},
		}

		require.Equal(t, zones, expected)
	})
}

type spyHttpClient struct {
	req  *http.Request
	resp string
}

func newSpyHttpClient() *spyHttpClient {
	return &spyHttpClient{}
}

func (c *spyHttpClient) Do(r *http.Request) (*http.Response, error) {
	c.req = r
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(c.resp))),
	}, nil
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
