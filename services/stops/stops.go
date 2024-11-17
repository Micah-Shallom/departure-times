package stops

import (
	"context"
	"fmt"

	"github.com/Micah-Shallom/departure-times/external/external_models"
	"github.com/Micah-Shallom/departure-times/external/requests"
	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/redis/go-redis/v9"
)

func GetRouteConfigurations(logger *utility.Logger, extReq requests.ExternalRequest, agency_tag string, route_tag string) (external_models.GetRoutesConfigResponse, error) {
	var (
		routeConfigResponse external_models.GetRoutesConfigResponse
	)

	data := map[string]string{
		"agency_tag": agency_tag,
		"route_tag":  route_tag,
	}

	response, err := extReq.SendExternalRequest(requests.GetRouteConfig, data)
	if err != nil {
		logger.Error("Error sending request: %s", err)
		return routeConfigResponse, err
	}

	routeConfigurations, ok := response.(external_models.GetRoutesConfigResponse)
	if !ok {
		logger.Error("Error casting response to GetRoutesConfigResponse")
		return routeConfigResponse, fmt.Errorf("error casting response to GetRoutesConfigResponse")
	}

	logger.Info("routeconfig fetched successfully")
	return routeConfigurations, nil
}

func GetStops(logger *utility.Logger, extReq requests.ExternalRequest, redisClient *redis.Client, reqdata map[string]interface{}) ([]external_models.StopInfo, error) {
	var (
		geoKey              = "geo:all_stops"
		ctx                 = context.Background()
		long                = reqdata["Longitude"].(string)
		lat                 = reqdata["Latitude"].(string)
		rad                 = reqdata["Radius"].(int)
		stopResultsResponse []external_models.StopInfo
	)

	stopResults, err := redisClient.GeoSearch(ctx, geoKey, &redis.GeoSearchQuery{
		Radius:    utility.ParseFloat(rad),
		Longitude: utility.ParseFloat(long),
		Latitude:  utility.ParseFloat(lat),
	}).Result()
	if err != nil {
		logger.Error("Error fetching geo search results: %s", err)
		return stopResultsResponse, err
	}

	for _, stopKey := range stopResults {
		stopData, err := redisClient.HGetAll(ctx, stopKey).Result()
		if err != nil {
			logger.Error("Error fetching stop data: %s", err)
			return stopResultsResponse, err
		}

		switch {
		case reqdata["AgencyTag"] == "" && reqdata["RouteTag"] == "":
			// Case 1: Both AgencyTag and RouteTag are empty (return all stops)
			utility.AppendStop(stopData, &stopResultsResponse)
	
		case reqdata["AgencyTag"] == "" && reqdata["RouteTag"] != "" && reqdata["RouteTag"] == stopData["routeTag"]:
			// Case 2: AgencyTag is empty, but RouteTag is specified and matches
			utility.AppendStop(stopData, &stopResultsResponse)
	
		case reqdata["RouteTag"] == "" && reqdata["AgencyTag"] != "" && reqdata["AgencyTag"] == stopData["agencyTag"]:
			// Case 3: RouteTag is empty, but AgencyTag is specified and matches
			utility.AppendStop(stopData, &stopResultsResponse)
	
		case reqdata["AgencyTag"] == stopData["agencyTag"] || reqdata["RouteTag"] == stopData["routeTag"]:
			// Case 4: Both AgencyTag and RouteTag are specified, and at least one matches
			utility.AppendStop(stopData, &stopResultsResponse)
		}
		

	}
	return stopResultsResponse, nil
}
