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

type RouteListResponse struct {
	Routes []RouteList `xml:"route"`
}

type RouteList struct {
	Tag   string `xml:"tag,attr"`
	Title string `xml:"title,attr"`
}

type GetRoutesResponse struct {
	Routes []Route `xml:"route"`
}

type GetRoutesConfigResponse struct {
	Route Route `xml:"route"`
}

type Route struct {
	Path          []Path      `xml:"path"`
	LonMax        string      `xml:"lonMax,attr"`
	Color         string      `xml:"color,attr"`
	OppositeColor string      `xml:"oppositeColor,attr"`
	Stop          []Stop      `xml:"stop"`
	Tag           string      `xml:"tag,attr"`
	LatMin        string      `xml:"latMin,attr"`
	Title         string      `xml:"title,attr"`
	LatMax        string      `xml:"latMax,attr"`
	LonMin        string      `xml:"lonMin,attr"`
	Direction     []Direction `xml:"direction"`
}

type Path struct {
	Point []Point `xml:"point"`
}

type Point struct {
	Lon string `xml:"lon,attr"`
	Lat string `xml:"lat,attr"`
}

type Stop struct {
	StopID     string `xml:"stopId,attr,omitempty"`
	Lon        string `xml:"lon,attr"`
	Tag        string `xml:"tag,attr"`
	ShortTitle string `xml:"shortTitle,attr,omitempty"`
	Title      string `xml:"title,attr"`
	Lat        string `xml:"lat,attr"`
}

type Direction struct {
	Stop     []DirectionStop `xml:"stop"`
	Name     string          `xml:"name,attr"`
	UseForUI string          `xml:"useForUI,attr"`
	Tag      string          `xml:"tag,attr"`
	Title    string          `xml:"title,attr"`
}

type DirectionStop struct {
	Tag string `xml:"tag,attr"`
}

type StopInfo struct {
	Tag        string  `json:"tag"`
	Title      string  `json:"title"`
	ShortTitle string  `json:"shortTitle"`
	Lat        string  `json:"lat"`
	Lon        string  `json:"lon"`
	StopID     string  `json:"stopId"`
	RouteTag   string  `json:"routeTag"`
	AgencyTag  string  `json:"agencyTag"`
	Distance   float64 `json:"distance"`
}

type GetStopsRequest struct {
	AgencyTag *string `json:"agencyTag"`
	RouteTag  *string `json:"routeTag"`
	Longitude string  `json:"longitude" validate:"required"`
	Latitude  string  `json:"latitude" validate:"required"`
	Radius    int     `json:"radius" validate:"required"`
}
