package caic

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type Client struct {
	http    doer
	caicURL string
}

type Zone struct {
	Name   string
	Url    string
	Rating int
}

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

func NewClient(caicURL string, http doer) *Client {
	return &Client{
		http:    http,
		caicURL: caicURL,
	}
}

func (c *Client) CanConnect() bool {
	resp, err := c.doRequest()
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		return false
	}
	resp.Body.Close()

	return true
}

func (c *Client) StateSummary() ([]Zone, error) {
	resp, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseResponse(string(b)), nil
}

func (c *Client) doRequest() (*http.Response, error) {
	req, err := http.NewRequest("GET", c.caicURL+"/caic/fx_map.php", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprint("unexpected status code ", resp.StatusCode))
	}
	return resp, nil
}

func parseResponse(caicResponse string) []Zone {
	r := `zone\[(\d+)\]='(.+)';\nurl\[\d+\]='(.+)';\nrating\[\d+\]=(\d)`
	regex := *regexp.MustCompile(r)
	matches := regex.FindAllStringSubmatch(caicResponse, -1)

	var z []Zone
	for _, m := range matches {
		z = append(z, Zone{
			// Index:  m[1], TODO: I don't need it right now
			Name:   m[2],
			Url:    m[3],
			Rating: toInt(m[4]),
		})
	}
	return z
}

func toInt(num string) int {
	n, _ := strconv.Atoi(num)
	return n
}
