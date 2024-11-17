package utility

import (
	"strconv"

	"github.com/Micah-Shallom/departure-times/external/external_models"
)

func ParseFloat(value any) float64 {
	switch v := value.(type) {
	case int:
		return float64(v)
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0
		}
		return f
	default:
		return 0
	}
}

func AppendStop(stopData map[string]string, stopResultsResponse *[]external_models.StopInfo) {
	*stopResultsResponse = append(*stopResultsResponse, external_models.StopInfo{
		Tag:        stopData["tag"],
		Title:      stopData["title"],
		ShortTitle: stopData["shortTitle"],
		Lat:        stopData["lat"],
		Lon:        stopData["lon"],
		StopID:     stopData["stopId"],
		RouteTag:   stopData["routeTag"],
		AgencyTag:  stopData["agencyTag"],
	})
}
