package caic

import "fmt"

type AspectDanger struct {
	Region        Region
	BelowTreeline OrdinalDanger
	NearTreeline  OrdinalDanger
	AboveTreeline OrdinalDanger
}

type OrdinalDanger struct {
	North     bool
	NorthEast bool
	East      bool
	SouthEast bool
	South     bool
	SouthWest bool
	West      bool
	NorthWest bool
}

var (
	ordinals   = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	elevations = []string{"Btl", "Tln", "Alp"}
)

func (c *Client) RegionAspectDanger(r Region) (AspectDanger, error) {
	resp, err := c.doRequest(fmt.Sprintf(regionPath, r))
	if err != nil {
		return AspectDanger{}, err
	}

	doc, err := toDocument(resp)
	if err != nil {
		return AspectDanger{}, err
	}

	aspects := make(map[int]map[int]bool)
	for ei, e := range elevations {
		od := make(map[int]bool)
		for oi, o := range ordinals {
			id := fmt.Sprintf("#%s%s_0.on", o, e)
			element := doc.Find(id)
			od[oi] = element.Nodes != nil // The element is on or it doesn't exit
		}
		aspects[ei] = od
	}

	return AspectDanger{
		Region: r,
		BelowTreeline: OrdinalDanger{
			North:     aspects[0][0],
			NorthEast: aspects[0][1],
			East:      aspects[0][2],
			SouthEast: aspects[0][3],
			South:     aspects[0][4],
			SouthWest: aspects[0][5],
			West:      aspects[0][6],
			NorthWest: aspects[0][7],
		},
		NearTreeline: OrdinalDanger{
			North:     aspects[1][0],
			NorthEast: aspects[1][1],
			East:      aspects[1][2],
			SouthEast: aspects[1][3],
			South:     aspects[1][4],
			SouthWest: aspects[1][5],
			West:      aspects[1][6],
			NorthWest: aspects[1][7],
		},
		AboveTreeline: OrdinalDanger{
			North:     aspects[2][0],
			NorthEast: aspects[2][1],
			East:      aspects[2][2],
			SouthEast: aspects[2][3],
			South:     aspects[2][4],
			SouthWest: aspects[2][5],
			West:      aspects[2][6],
			NorthWest: aspects[2][7],
		},
	}, nil
}
