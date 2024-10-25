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
		return agencyResponse, err
	}

	agencyListResponse, ok := response.(external_models.GetAgenciesResponse)
	if !ok {
		return agencyResponse, fmt.Errorf("error casting response to GetAgenciesResponse")
	}

	return agencyListResponse, nil
}

func GetRouteList(logger *utility.Logger, extReq requests.ExternalRequest, agency_tag string) (external_models.GetRoutesResponse, error) {
	var routeResponse external_models.GetRoutesResponse

	response, err := extReq.SendExternalRequest(requests.GetRouteList, agency_tag)
	if err != nil {
		return routeResponse, err
	}

	routeListResponse, ok := response.(external_models.GetRoutesResponse)
	if !ok {
		return routeResponse, fmt.Errorf("error casting response to GetAgenciesResponse")
	}

	return routeListResponse, nil
}
