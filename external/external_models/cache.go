package external_models


type Cache struct {
	RouteListCache map[string]GetRoutesResponse // map[agency_tag]GetRoutesResponse
	RouteConfigCache map[string]GetRoutesConfigResponse // map[route_tag]GetRoutesConfigResponse
	StopListCache map[string]map[string][]Stop 
}