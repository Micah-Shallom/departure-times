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

func (base *Controller) GetRouteConfigList(c *gin.Context) {
	var (
		agency_tag = c.Query("agency_tag")
		route_tag  = c.Query("route_tag")
	)

	if agency_tag == "" || route_tag == "" {
		base.Logger.Error("agency tag or route tag is empty. Provide an agency tag and route tag")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "agency tag or route tag is empty, provide an agency tag and route tag.", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	response, err := service.GetRouteConfigurations(base.Logger, base.ExtReq, agency_tag, route_tag)
	if err != nil {
		base.Logger.Error("unable to get route configurations", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get route configurations", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("route configurations fetched successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "route configurations fetched successfully", response)
	c.JSON(http.StatusOK, rd)
}
