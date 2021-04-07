package caic

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type CaicClient struct {
	http    doer
	caicURL string
}

type Zone struct {
	Index  string
	Name   string
	Url    string
	Rating string
}

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

func NewClient(caicURL string, http doer) *CaicClient {
	return &CaicClient{
		http:    http,
		caicURL: caicURL,
	}
}

func (c *CaicClient) StateSummary() []Zone {
	req, err := http.NewRequest("get", c.caicURL+"/caic/fx_map.php", nil)
	if err != nil {
		//TODO, how to propagate errors in the plugin
		log.Printf("error creating request")
		return nil
	}

	resp, err := c.http.Do(req)
	if err != nil {
		//TODO, how to propagate errors in the plugin
		log.Printf("error making request to website")
		return nil
	}
	defer resp.Body.Close()

	//TODO what about an error from the wobsite?
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//TODO, how to propagate errors in the plugin
		log.Printf("error reading response")
		return nil
	}

	return parseResponse(string(b))
}

func parseResponse(caicResponse string) []Zone {
	r := `zone\[(\d+)\]='(.+)';\nurl\[\d+\]='(.+)';\nrating\[\d+\]=(\d)`
	regex := *regexp.MustCompile(r)
	matches := regex.FindAllStringSubmatch(caicResponse, -1)

	var z []Zone
	for _, m := range matches {
		z = append(z, Zone{
			Index:  m[1],
			Name:   m[2],
			Url:    m[3],
			Rating: m[4],
		})
	}
	return z
}
