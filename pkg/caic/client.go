package caic

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	http    doer
	caicURL string
}

type Zone struct {
	Index         Region
	Name          string
	Rating        int
	AboveTreeline int
	NearTreeline  int
	BelowTreeline int
}

type elevation int

const (
	aboveTreeline = iota + 1
	nearTreeline
	belowTreeline
)

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

func NewClient(caicURL string, http doer) *Client {
	return &Client{
		http:    http,
		caicURL: caicURL,
	}
}

func (c *Client) StateSummary() ([]Zone, error) {
	resp, err := c.doRequest("/caic/fx_map.php")
	if err != nil {
		return nil, err
	}

	zoneIdRatingPattern := `zone\[(\d+)\]='(.+)';\nurl\[\d+\]='(.+)';\nrating\[\d+\]=(\d)`
	regex := *regexp.MustCompile(zoneIdRatingPattern)
	matches := regex.FindAllStringSubmatch(resp, -1)

	var z []Zone
	for _, m := range matches {
		z = append(z, Zone{
			Index:  Region(toInt(m[1])),
			Name:   m[2],
			Rating: toInt(m[4]),
		})
	}

	return z, nil
}

func (c *Client) RegionSummary(r Region) (Zone, error) {
	path := fmt.Sprintf("/caic/pub_bc_avo.php?zone_id=%d", r)
	resp, err := c.doRequest(path)
	if err != nil {
		return Zone{}, err
	}

	doc, err := toDocument(resp)
	if err != nil {
		return Zone{}, err
	}

	z := Zone{
		Index:         r,
		Name:          r.String(),
		AboveTreeline: ratingFor(aboveTreeline, doc),
		NearTreeline:  ratingFor(nearTreeline, doc),
		BelowTreeline: ratingFor(belowTreeline, doc),
	}
	z.Rating = max(z.AboveTreeline, z.NearTreeline, z.BelowTreeline)

	return z, nil
}

func (c *Client) CanConnect() bool {
	_, err := c.doRequest("/caic/fx_map.php")
	return err == nil
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

func toDocument(s string) (*goquery.Document, error) {
	r := strings.NewReader(s)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func ratingFor(e elevation, doc *goquery.Document) int {
	query := fmt.Sprintf("#avalanche-forecast > table.table.table-striped-body.table-treeline > tbody > tr:nth-child(%d) > td.today-text > strong", e)
	ratingText := doc.Find(query).Nodes[0].FirstChild.Data

	return parseRating(ratingText)
}

func parseRating(s string) int {
	ratingPattern := `.+\((\d)\)`
	regex := *regexp.MustCompile(ratingPattern)
	matches := regex.FindAllStringSubmatch(s, -1)

	return toInt(matches[0][1])
}

func toInt(num string) int {
	n, _ := strconv.Atoi(num)
	return n
}

func max(i ...int) int {
	m := -1
	for _, n := range i {
		if n > m {
			m = n
		}
	}
	return m
}
