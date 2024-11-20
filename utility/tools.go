package utility

import (
	"math"
	"sort"
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

func AppendStop(stopData map[string]string, stopResultsResponse *[]external_models.StopInfo, distance float64) {
	*stopResultsResponse = append(*stopResultsResponse, external_models.StopInfo{
		Tag:        stopData["tag"],
		Title:      stopData["title"],
		ShortTitle: stopData["shortTitle"],
		Lat:        stopData["lat"],
		Lon:        stopData["lon"],
		StopID:     stopData["stopId"],
		RouteTag:   stopData["routeTag"],
		AgencyTag:  stopData["agencyTag"],
		Distance:   distance,
	})
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 //  the radius of the earth in kilometers
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

func SortStopDistance(stopResultsResponse []external_models.StopInfo) {
	sort.Slice(stopResultsResponse, func(i, j int) bool {
		return stopResultsResponse[i].Distance < stopResultsResponse[j].Distance
	})
}