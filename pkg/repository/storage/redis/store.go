package redis

import (
	"fmt"
	"strconv"

	"github.com/Micah-Shallom/departure-times/external/external_models"
	"github.com/redis/go-redis/v9"
)

func storeStopsInRedis(redisClient *redis.Client, agencyTag, routeTag string, stops []external_models.Stop) error {
	geoKey := "geo:all_stops"

	parseFloat := func(value string) float64 {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0
		}
		return f
	}


	for _, stop := range stops {

		err := redisClient.GeoAdd(ctx, geoKey, &redis.GeoLocation{
			Longitude: parseFloat(stop.Lon),
			Latitude:  parseFloat(stop.Lat),
			Name:      fmt.Sprintf("stop:%s", stop.Tag),
		}).Err()
		if err != nil {
			return fmt.Errorf("error storing stops in redis: %s", err)
		}

		stopKey := fmt.Sprintf("stop:%s", stop.Tag)
		_, err = redisClient.HMSet(ctx, stopKey, map[string]any{
			"tag":        stop.Tag,
			"title":      stop.Title,
			"shortTitle": stop.ShortTitle,
			"lat":        stop.Lat,
			"lon":        stop.Lon,
			"stopId":     stop.StopID,
			"routeTag":   routeTag,
			"agencyTag":  agencyTag,
		}).Result()
		if err != nil {
			return fmt.Errorf("error storing stops in redis hashset: %s", err)
		}
	}

	return nil
}
