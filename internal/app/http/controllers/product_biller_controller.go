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

// ProductBillerController defines the HTTP layer for ProductBiller entities.
type ProductBillerController struct {
	usecases usecases.ProductBillerUseCase
	logger   *zerolog.Logger
}

// NewProductBillerController creates a new instance of ProductBillerController.
func NewProductBillerController(usecases usecases.ProductBillerUseCase, logger *zerolog.Logger) *ProductBillerController {
	return &ProductBillerController{
		usecases: usecases,
		logger:   logger,
	}
}

const eventClassProductBiller = "controller.productBiller"

// Create handles POST requests to create a new ProductBiller.
func (c *ProductBillerController) Create(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	var productBiller *models.CreateProductBillerRequest
	if err := ctx.Bind(productBiller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(productBiller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.usecases.Create(reqCtx, productBiller.ToEntity()); err != nil {
		logger.Error(reqCtx, eventClassProductBiller, "Create", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"message": "Product Biller created successfully"})
}

// Update handles PUT requests to update an existing ProductBiller.
func (c *ProductBillerController) Update(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	var productBiller *models.UpdateProductBillerRequest
	if err := ctx.Bind(productBiller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(productBiller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.usecases.Update(reqCtx, uint(id), productBiller.ToEntity()); err != nil {
		logger.Error(reqCtx, eventClassProductBiller, "Update", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Product Biller updated successfully"})
}

// Delete handles DELETE requests to remove a ProductBiller.
func (c *ProductBillerController) Delete(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	if err := c.usecases.Delete(reqCtx, uint(id)); err != nil {
		logger.Error(reqCtx, eventClassProductBiller, "Delete", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Product Biller deleted successfully"})
}

// FetchOne handles GET requests to retrieve a single ProductBiller.
func (c *ProductBillerController) FetchOne(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	productBiller, err := c.usecases.FetchOne(reqCtx, uint(id))
	if err != nil {
		logger.Error(reqCtx, eventClassProductBiller, "FetchOne", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, productBiller.ToResponse())
}

// FetchMany handles GET requests to retrieve multiple ProductBillers based on filters.
func (c *ProductBillerController) FetchMany(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	filter := make(map[string]interface{})
	if productID := ctx.QueryParam("product_id"); productID != "" {
		filter["product_id"] = productID
	}
	if billerID := ctx.QueryParam("biller_id"); billerID != "" {
		filter["biller_id"] = billerID
	}

	productBillers, err := c.usecases.FetchMany(reqCtx, filter)
	if err != nil {
		logger.Error(reqCtx, eventClassProductBiller, "FetchMany", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := utils.TransformSlice(productBillers, func(pb *models.ProductBiller) *models.ProductBillerResponse {
		return pb.ToResponse()
	})
	return ctx.JSON(http.StatusOK, response)
}

// FetchManyWithPagination handles GET requests to retrieve paginated ProductBillers.
func (c *ProductBillerController) FetchManyWithPagination(ctx echo.Context) error {
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
	if productID := ctx.QueryParam("product_id"); productID != "" {
		filter["product_id"] = productID
	}
	if billerID := ctx.QueryParam("biller_id"); billerID != "" {
		filter["biller_id"] = billerID
	}

	productBillers, pagination, err := c.usecases.FetchManyWithPagination(reqCtx, filter, page, limit)
	if err != nil {
		logger.Error(reqCtx, eventClassProductBiller, "FetchManyWithPagination", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := utils.TransformSlice(productBillers, func(pb *models.ProductBiller) *models.ProductBillerResponse {
		return pb.ToResponse()
	})
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data":       response,
		"pagination": pagination,
	})
}
