package service

import (
	"fmt"

	"github.com/Micah-Shallom/departure-times/external/external_models"
	"github.com/Micah-Shallom/departure-times/external/requests"
	"github.com/Micah-Shallom/departure-times/utility"
)

func GetAgencyList(logger *utility.Logger, extReq requests.ExternalRequest) (external_models.GetAgenciesResponse, error) {
	var agencyResponse external_models.GetAgenciesResponse

	response, err := extReq.SendExternalRequest(requests.GetAgencyList, nil)
	if err != nil {
		logger.Error("Error sending request: %s", err)
		return agencyResponse, err
	}

	agencyListResponse, ok := response.(external_models.GetAgenciesResponse)
	if !ok {
		logger.Error("Error casting response to GetAgenciesResponse")
		return agencyResponse, fmt.Errorf("error casting response to GetAgenciesResponse")
	}

	logger.Info("agency list fetched successfully")
	return agencyListResponse, nil
}

func GetRouteList(logger *utility.Logger, extReq requests.ExternalRequest, agency_tag string) (external_models.RouteListResponse, error) {
	var routeResponse external_models.RouteListResponse

	response, err := extReq.SendExternalRequest(requests.GetRouteList, agency_tag)
	if err != nil {
		logger.Error("Error sending request: %s", err)
		return routeResponse, err
	}

	routeListResponse, ok := response.(external_models.RouteListResponse)
	if !ok {
		logger.Error("Error casting response to GetAgenciesResponse")
		return routeResponse, fmt.Errorf("error casting response to GetAgenciesResponse")
	}

	logger.Info("route list fetched successfully")
	return routeListResponse, nil
}
