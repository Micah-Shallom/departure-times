package agency

import (
	"net/http"

	"github.com/Micah-Shallom/departure-times/external/requests"
	"github.com/Micah-Shallom/departure-times/internal/config"
	service "github.com/Micah-Shallom/departure-times/services/agency"
	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Logger *utility.Logger
	ExtReq requests.ExternalRequest
	Config *config.Config
}

func (base *Controller) GetAgencies(c *gin.Context) {

	response, err := service.GetAgencyList(base.Logger, base.ExtReq)
	if err != nil {
		base.Logger.Error("unable to fetch agency list", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to fetch agency list", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("agency list fetched successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "agency list fetched successfully", response)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetRouteList(c *gin.Context) {
	var (
		agency_tag = c.Query("agency_tag")
	)

	if agency_tag == "" {
		base.Logger.Error("agency tag is empty. Provide an agency tag")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "agency tag is empty, provide an agency tag.", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	response, err := service.GetRouteList(base.Logger, base.ExtReq, agency_tag)
	if err != nil {
		base.Logger.Error("unable to get routelist", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get routelist", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("route list fetched successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "route list fetched successfully", response)
	c.JSON(http.StatusOK, rd)
}
