package controller

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"golang-boilerplate/internal/pkg/models"
	service "golang-boilerplate/internal/pkg/services"
	"golang-boilerplate/internal/pkg/utils/common"
)

// BillerController defines the HTTP layer for Biller entities.
type BillerController struct {
	service service.BillerService
}

// NewBillerController creates a new instance of BillerController.
func NewBillerController(service service.BillerService) *BillerController {
	return &BillerController{
		service: service,
	}
}

// Create handles POST requests to create a new Biller.
func (c *BillerController) Create(ctx echo.Context) error {
	var biller *models.CreateBillerRequest
	if err := ctx.Bind(biller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(biller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.service.Create(ctx.Request().Context(), biller.ToEntity()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"message": "Biller created successfully"})
}

// Update handles PUT requests to update an existing Biller.
func (c *BillerController) Update(ctx echo.Context) error {
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

	if err := c.service.Update(ctx.Request().Context(), uint(id), biller.ToEntity()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Biller updated successfully"})
}

// Delete handles DELETE requests to remove a Biller.
func (c *BillerController) Delete(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	if err := c.service.Delete(ctx.Request().Context(), uint(id)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Biller deleted successfully"})
}

// GetOne handles GET requests to retrieve a single Biller.
func (c *BillerController) GetOne(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	biller, err := c.service.GetOne(ctx.Request().Context(), uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, biller.ToResponse())
}

// GetMany handles GET requests to retrieve multiple Billers based on filters.
func (c *BillerController) GetMany(ctx echo.Context) error {
	filter := make(map[string]interface{})
	if label := ctx.QueryParam("label"); label != "" {
		filter["label"] = label
	}

	billers, err := c.service.GetMany(ctx.Request().Context(), filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := common.TransformSlice(billers, func(biller *models.Biller) *models.BillerResponse {
		return biller.ToResponse()
	})
	return ctx.JSON(http.StatusOK, response)
}

// GetManyWithPagination handles GET requests to retrieve paginated Billers.
func (c *BillerController) GetManyWithPagination(ctx echo.Context) error {
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

	billers, pagination, err := c.service.GetManyWithPagination(ctx.Request().Context(), filter, page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := common.TransformSlice(billers, func(biller *models.Biller) *models.BillerResponse {
		return biller.ToResponse()
	})
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data":       response,
		"pagination": pagination,
	})
}
