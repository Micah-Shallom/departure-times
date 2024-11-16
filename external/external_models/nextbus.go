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

// custom unmarshaller for directions struct to handle slice and object responses
// func (d *DirectionOrDirections) UnmarshalJSON(data []byte) error {
// 	var (
// 		single   Direction
// 		multiple []Direction
// 	)

// 	if err := json.Unmarshal(data, &single); err == nil {
// 		d.Direction = []Direction{single}
// 		return nil
// 	}

// 	if err := json.Unmarshal(data, &multiple); err == nil {
// 		d.Direction = multiple
// 		return nil
// 	}
// 	return fmt.Errorf("failed to custom unmarshal direction")
// }