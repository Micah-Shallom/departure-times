package stops

import (
	"net/http"

	"github.com/Micah-Shallom/departure-times/external/requests"
	"github.com/Micah-Shallom/departure-times/internal/config"
	service "github.com/Micah-Shallom/departure-times/services/stops"
	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Logger *utility.Logger
	Config *config.Config
	ExtReq requests.ExternalRequest
}

func (base *Controller) GetStops(c *gin.Context) {
	var (
		agency_tag = c.Query("agency_tag")
		route_tag  = c.Query("route_tag")
	)

	response, err := service.GetRouteConfigurations(base.Logger, base.ExtReq, agency_tag, route_tag)
	if err != nil {
		base.Logger.Error("unable to get stops", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get stops", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("stops fetched successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "stops fetched successfully", response)
	c.JSON(http.StatusOK, rd)
}
