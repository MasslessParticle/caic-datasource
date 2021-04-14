package caic

type Region int

const (
	EntireState Region = iota - 1
	SteamboatFlatTops
	FrontRange
	VailSummitCounty
	SawatchRange
	Aspen
	Gunnison
	GrandMesa
	NorthernSanJuan
	SouthernSanJuan
	SangreDeCristo
)

func (d Region) String() string {
	if d == EntireState {
		return "Entire State"
	}

	return []string{
		"Steamboat & Flat Tops",
		"Front Range",
		"Vail & Summit County",
		"Sawatch Range",
		"Aspen",
		"Gunnison",
		"Grand Mesa",
		"Northern San Juan",
		"Southern San Juan",
		"Sangre de Cristo",
	}[d]
}
