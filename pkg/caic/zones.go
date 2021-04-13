package caic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

func (c *Client) RegionSummary(r Region) ([]Zone, error) {
	if r == EntireState {
		return c.stateSummary()
	}
	return c.singleRegionSummary(r)
}

func (c *Client) CanConnect() bool {
	_, err := c.doRequest(homePath)
	return err == nil
}

func (c *Client) stateSummary() ([]Zone, error) {
	var zones []Zone
	for i := SteamboatFlatTops; i <= SangreDeCristo; i++ {
		z, err := c.singleRegionSummary(Region(i))
		if err != nil {
			return nil, err
		}
		zones = append(zones, z...)
	}

	return zones, nil
}

func (c *Client) singleRegionSummary(r Region) ([]Zone, error) {
	path := fmt.Sprintf(regionPath, r)
	resp, err := c.doRequest(path)
	if err != nil {
		return nil, err
	}

	doc, err := toDocument(resp)
	if err != nil {
		return nil, err
	}

	z := Zone{
		Index:         r,
		Name:          r.String(),
		AboveTreeline: ratingFor(aboveTreeline, doc),
		NearTreeline:  ratingFor(nearTreeline, doc),
		BelowTreeline: ratingFor(belowTreeline, doc),
	}
	z.Rating = max(z.AboveTreeline, z.NearTreeline, z.BelowTreeline)

	return []Zone{z}, nil
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

	if len(matches) > 0 {
		return toInt(matches[0][1])
	}
	return 0
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
