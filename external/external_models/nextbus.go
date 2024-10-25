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
	Routes []struct {
		Tag           string  `xml:"tag,attr"`
		Title         string  `xml:"title,attr"`
		Color         string  `xml:"color,attr"`
		OppositeColor string  `xml:"oppositeColor,attr"`
		LatMin        float64 `xml:"latMin,attr"`
		LatMax        float64 `xml:"latMax,attr"`
		LonMin        float64 `xml:"lonMin,attr"`
		LonMax        float64 `xml:"lonMax,attr"`

		Stops []struct {
			Tag        string  `xml:"tag,attr"`
			Title      string  `xml:"title,attr"`
			ShortTitle string  `xml:"shortTitle,attr,omitempty"`
			Lat        float64 `xml:"lat,attr"`
			Lon        float64 `xml:"lon,attr"`
			StopID     string  `xml:"stopId,attr,omitempty"`
		} `xml:"stop"`

		Directions []struct {
			Tag      string   `xml:"tag,attr"`
			Title    string   `xml:"title,attr"`
			Name     string   `xml:"name,attr,omitempty"`
			UseForUI bool     `xml:"useForUI,attr"`
			StopTags []string `xml:"stop"`
		} `xml:"direction"`

		Paths []struct {
			Points []struct {
				Lat float64 `xml:"lat,attr"`
				Lon float64 `xml:"lon,attr"`
			} `xml:"point"`
		} `xml:"path"`
	} `xml:"route"`
}
