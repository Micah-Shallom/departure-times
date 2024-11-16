package router

import (
	"fmt"

	"github.com/Micah-Shallom/departure-times/external/requests"
	"github.com/Micah-Shallom/departure-times/internal/config"
	"github.com/Micah-Shallom/departure-times/pkg/controller/agency"
	"github.com/Micah-Shallom/departure-times/pkg/controller/stops"
	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/gin-gonic/gin"
)

func NextBus(r *gin.Engine, ApiVersion string, logger *utility.Logger, cfg *config.Config) {
	extReq := requests.ExternalRequest{Logger: logger}
	nextbusURL := r.Group(fmt.Sprintf("%s/nextbus", ApiVersion))
	agency := agency.Controller{Logger: logger, ExtReq: extReq, Config: cfg}
	stops := stops.Controller{Logger: logger, ExtReq: extReq, Config: cfg}

	nextbusURL.GET("/agency", agency.GetAgencies)
	nextbusURL.GET("/routelist", agency.GetRouteList)

	nextbusURL.GET("/routeconfig", stops.GetRouteConfigList)
}
