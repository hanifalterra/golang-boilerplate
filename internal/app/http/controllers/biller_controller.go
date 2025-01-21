package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"golang-boilerplate/internal/app/http/usecases"
	"golang-boilerplate/internal/pkg/logger"
	"golang-boilerplate/internal/pkg/models"
	"golang-boilerplate/internal/pkg/utils"
)

// BillerController defines the HTTP layer for Biller entities.
type BillerController struct {
	usecases usecases.BillerUseCase
	logger   *zerolog.Logger
}

// NewBillerController creates a new instance of BillerController.
func NewBillerController(usecases usecases.BillerUseCase, logger *zerolog.Logger) *BillerController {
	return &BillerController{
		usecases: usecases,
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

	if err := c.usecases.Create(reqCtx, biller.ToEntity()); err != nil {
		logger.Error(reqCtx, eventClassBiller, "Create", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"message": "Biller created successfully"})
}

// Update handles PUT requests to update an existing Biller.
func (c *BillerController) Update(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.Atoi(ctx.Param("id"))
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

	if err := c.usecases.Update(reqCtx, id, biller.ToEntity()); err != nil {
		logger.Error(reqCtx, eventClassBiller, "Update", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Biller updated successfully"})
}

// Delete handles DELETE requests to remove a Biller.
func (c *BillerController) Delete(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	if err := c.usecases.Delete(reqCtx, id); err != nil {
		logger.Error(reqCtx, eventClassBiller, "Delete", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Biller deleted successfully"})
}

// FetchOne handles GET requests to retrieve a single Biller.
func (c *BillerController) FetchOne(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	biller, err := c.usecases.FetchOne(reqCtx, id)
	if err != nil {
		logger.Error(reqCtx, eventClassBiller, "FetchOne", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, biller.ToResponse())
}

// FetchMany handles GET requests to retrieve multiple Billers based on filters.
func (c *BillerController) FetchMany(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	filter := make(map[string]interface{})
	if label := ctx.QueryParam("label"); label != "" {
		filter["label"] = label
	}

	billers, err := c.usecases.FetchMany(reqCtx, filter)
	if err != nil {
		logger.Error(reqCtx, eventClassBiller, "FetchMany", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := utils.TransformSlice(billers, func(biller *models.Biller) *models.BillerResponse {
		return biller.ToResponse()
	})
	return ctx.JSON(http.StatusOK, response)
}

// FetchManyWithPagination handles GET requests to retrieve paginated Billers.
func (c *BillerController) FetchManyWithPagination(ctx echo.Context) error {
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

	billers, pagination, err := c.usecases.FetchManyWithPagination(reqCtx, filter, page, limit)
	if err != nil {
		logger.Error(reqCtx, eventClassBiller, "FetchManyWithPagination", err.Error())
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
