package external_models

type GetAgenciesResponse struct {
	Agencies []Agency `xml:"agency"`
}

type Agency struct {
	Tag         string `xml:"tag,attr"`
	Title       string `xml:"title,attr"`
	ShortTitle  string `xml:"shortTitle,attr"`
	RegionTitle string `xml:"regionTitle,attr"`
}

type GetRoutesResponse struct {
	Routes []struct {
		Tag   string `xml:"tag,attr"`
		Title string `xml:"title,attr"`
	} `xml:"route"`
}

type GetRoutesConfigResponse struct {
	Config RouteWrapper `xml:"route"`
}

type RouteWrapper struct {
	Route      Route       `xml:",inline"`
	Stops      []Stop      `xml:"stop"`
	Directions []Direction `xml:"direction"`
	Paths      []Path      `xml:"path"`
}

type Route struct {
	Tag           string  `xml:"tag,attr"`
	Title         string  `xml:"title,attr"`
	Color         string  `xml:"color,attr"`
	OppositeColor string  `xml:"oppositeColor,attr"`
	LatMin        float64 `xml:"latMin,attr"`
	LatMax        float64 `xml:"latMax,attr"`
	LonMin        float64 `xml:"lonMin,attr"`
	LonMax        float64 `xml:"lonMax,attr"`
}

type Stop struct {
	Tag        string  `xml:"tag,attr"`
	Title      string  `xml:"title,attr"`
	ShortTitle string  `xml:"shortTitle,attr,omitempty"`
	Lat        float64 `xml:"lat,attr"`
	Lon        float64 `xml:"lon,attr"`
	StopID     string  `xml:"stopId,attr,omitempty"`
}

type Direction struct {
	Tag      string   `xml:"tag,attr"`
	Title    string   `xml:"title,attr"`
	Name     string   `xml:"name,attr,omitempty"`
	UseForUI bool     `xml:"useForUI,attr"`
	StopTags []string `xml:"stop>tag"`
}

type Path struct {
	Points []Point `xml:"point"`
}

type Point struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
}
