package requests

import (
	"fmt"

	"github.com/Micah-Shallom/departure-times/external/third_party/nextbus"
	"github.com/Micah-Shallom/departure-times/internal/config"
	"github.com/Micah-Shallom/departure-times/utility"
)

type ExternalRequest struct {
	Logger *utility.Logger
}

var (
	JSONDecodeMethod string = "json"
	XMLDecodeMethod  string = "xml"
	GetAgencyList    string = "get_agency_list"
	GetRouteList     string = "get_route_list"
)

func (er *ExternalRequest) SendExternalRequest(name string, data interface{}) (interface{}, error) {
	var (
		config = config.GetConfig(er.Logger)
	)

	switch name {
	case GetAgencyList:
		obj := nextbus.RequestObj{
			Name:         name,
			Path:         config.App.NextBusURL,
			Method:       "GET",
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: XMLDecodeMethod,
			Logger:       er.Logger,
		}
		return obj.GetAgencyList()
	case GetRouteList:
		data := data.(string)

		obj := nextbus.RequestObj{
			Name:         name,
			Path:         config.App.NextBusURL,
			Method:       "GET",
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: XMLDecodeMethod,
			Logger:       er.Logger,
		}
		return obj.GetRouteList(data)
	default:
		return nil, fmt.Errorf("request name %s not found", name)
	}
}
