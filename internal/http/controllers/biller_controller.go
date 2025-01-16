package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"golang-boilerplate/internal/http/services"
	"golang-boilerplate/internal/pkg/logger"
	"golang-boilerplate/internal/pkg/models"
	"golang-boilerplate/internal/pkg/utils"
)

// BillerController defines the HTTP layer for Biller entities.
type BillerController struct {
	services services.BillerService
	logger   *zerolog.Logger
}

// NewBillerController creates a new instance of BillerController.
func NewBillerController(services services.BillerService, logger *zerolog.Logger) *BillerController {
	return &BillerController{
		services: services,
		logger:   logger,
	}
}

const eventClassBiller = "controller.biller"

// Create handles POST requests to create a new Biller.
func (c *BillerController) Create(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	var biller *models.CreateBillerRequest
	if err := ctx.Bind(biller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(biller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.services.Create(reqCtx, biller.ToEntity()); err != nil {
		logger.Error(reqCtx, eventClassBiller, "Create", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"message": "Biller created successfully"})
}

// Update handles PUT requests to update an existing Biller.
func (c *BillerController) Update(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	var biller *models.UpdateBillerRequest
	if err := ctx.Bind(biller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(biller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.services.Update(reqCtx, uint(id), biller.ToEntity()); err != nil {
		logger.Error(reqCtx, eventClassBiller, "Update", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Biller updated successfully"})
}

// Delete handles DELETE requests to remove a Biller.
func (c *BillerController) Delete(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	if err := c.services.Delete(reqCtx, uint(id)); err != nil {
		logger.Error(reqCtx, eventClassBiller, "Delete", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Biller deleted successfully"})
}

// GetOne handles GET requests to retrieve a single Biller.
func (c *BillerController) GetOne(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	biller, err := c.services.GetOne(reqCtx, uint(id))
	if err != nil {
		logger.Error(reqCtx, eventClassBiller, "GetOne", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, biller.ToResponse())
}

// GetMany handles GET requests to retrieve multiple Billers based on filters.
func (c *BillerController) GetMany(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	filter := make(map[string]interface{})
	if label := ctx.QueryParam("label"); label != "" {
		filter["label"] = label
	}

	billers, err := c.services.GetMany(reqCtx, filter)
	if err != nil {
		logger.Error(reqCtx, eventClassBiller, "GetMany", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := utils.TransformSlice(billers, func(biller *models.Biller) *models.BillerResponse {
		return biller.ToResponse()
	})
	return ctx.JSON(http.StatusOK, response)
}

// GetManyWithPagination handles GET requests to retrieve paginated Billers.
func (c *BillerController) GetManyWithPagination(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	filter := make(map[string]interface{})
	if label := ctx.QueryParam("label"); label != "" {
		filter["label"] = label
	}

	billers, pagination, err := c.services.GetManyWithPagination(reqCtx, filter, page, limit)
	if err != nil {
		logger.Error(reqCtx, eventClassBiller, "GetManyWithPagination", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := utils.TransformSlice(billers, func(biller *models.Biller) *models.BillerResponse {
		return biller.ToResponse()
	})
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data":       response,
		"pagination": pagination,
	})
}
