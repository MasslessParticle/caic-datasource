package caic

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	http    doer
	caicURL string
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

func (c *Client) doRequest(path string) (string, error) {
	req, err := http.NewRequest("GET", c.caicURL+path, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprint("unexpected status code ", resp.StatusCode))
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
