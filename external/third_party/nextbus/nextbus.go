package nextbus

import (
	"fmt"

	"github.com/Micah-Shallom/departure-times/external/external_models"
)

func (r *RequestObj) GetAgencyList() (external_models.GetAgenciesResponse, error) {
	var (
		logger           = r.Logger
		outBoundResponse = external_models.GetAgenciesResponse{}
	)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	path := "?command=agencyList"

	err := r.getNewSendRequestObject(nil, headers, path).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("Error sending request: %s", err)
		return outBoundResponse, err
	}

	return outBoundResponse, nil
}

func (r *RequestObj) GetRouteList(agency string) (external_models.GetRoutesResponse, error) {
	var (
		logger           = r.Logger
		outBoundResponse = external_models.GetRoutesResponse{}
	)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	path := fmt.Sprintf("?command=routeList&a=%s", agency)

	err := r.getNewSendRequestObject(nil, headers, path).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("Error sending request: %s", err)
		return outBoundResponse, err
	}

	return outBoundResponse, nil
}
