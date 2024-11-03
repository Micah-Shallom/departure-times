package stops

import (
	"fmt"

	"github.com/Micah-Shallom/departure-times/external/external_models"
	"github.com/Micah-Shallom/departure-times/external/requests"
	"github.com/Micah-Shallom/departure-times/utility"
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
