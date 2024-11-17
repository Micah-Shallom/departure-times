package stops

import (
	"net/http"

	"github.com/Micah-Shallom/departure-times/external/external_models"
	"github.com/Micah-Shallom/departure-times/external/requests"
	"github.com/Micah-Shallom/departure-times/internal/config"
	"github.com/Micah-Shallom/departure-times/pkg/repository/storage"
	service "github.com/Micah-Shallom/departure-times/services/stops"
	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Controller struct {
	Logger    *utility.Logger
	Config    *config.Config
	ExtReq    requests.ExternalRequest
	Validator *validator.Validate
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

func (base *Controller) GetStops(c *gin.Context) {
	var (
		agencyTag                 = c.Query("agency")
		routeTag                  = c.Query("route")
		redisClient *redis.Client = storage.DB.Redis
		req                       = external_models.GetStopsRequest{}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		base.Logger.Error("unable to bind request", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to bind request", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err := base.Validator.Struct(&req)
	if err != nil {
		base.Logger.Info("Validation failed", err)
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	data := map[string]interface{}{
		"Longitude": req.Longitude,
		"Latitude":  req.Latitude,
		"Radius":    req.Radius,
		"Agency":    agencyTag,
		"Route":     routeTag,
	}

	response, err := service.GetStops(base.Logger, base.ExtReq, redisClient, data)
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
