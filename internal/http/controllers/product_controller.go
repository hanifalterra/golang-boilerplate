package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"golang-boilerplate/internal/pkg/logger"
	"golang-boilerplate/internal/pkg/models"
	"golang-boilerplate/internal/pkg/services"
	"golang-boilerplate/internal/pkg/utils"
)

// ProductController defines the HTTP layer for Product entities.
type ProductController struct {
	services services.ProductService
	logger   *zerolog.Logger
}

// NewProductController creates a new instance of ProductController.
func NewProductController(services services.ProductService, logger *zerolog.Logger) *ProductController {
	return &ProductController{
		services: services,
		logger:   logger,
	}
}

const eventClassProduct = "controller.product"

// Create handles POST requests to create a new Product.
func (c *ProductController) Create(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	var product *models.CreateProductRequest
	if err := ctx.Bind(product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.services.Create(reqCtx, product.ToEntity()); err != nil {
		logger.Error(reqCtx, eventClassProduct, "Create", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"message": "Product created successfully"})
}

// Update handles PUT requests to update an existing Product.
func (c *ProductController) Update(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	var product *models.UpdateProductRequest
	if err := ctx.Bind(product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.services.Update(reqCtx, uint(id), product.ToEntity()); err != nil {
		logger.Error(reqCtx, eventClassProduct, "Update", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Product updated successfully"})
}

// Delete handles DELETE requests to remove a Product.
func (c *ProductController) Delete(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	if err := c.services.Delete(reqCtx, uint(id)); err != nil {
		logger.Error(reqCtx, eventClassProduct, "Delete", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

// GetOne handles GET requests to retrieve a single Product.
func (c *ProductController) GetOne(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	product, err := c.services.GetOne(reqCtx, uint(id))
	if err != nil {
		logger.Error(reqCtx, eventClassProduct, "GetOne", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, product.ToResponse())
}

// GetMany handles GET requests to retrieve multiple Products based on filters.
func (c *ProductController) GetMany(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	filter := make(map[string]interface{})
	if label := ctx.QueryParam("label"); label != "" {
		filter["label"] = label
	}

	products, err := c.services.GetMany(reqCtx, filter)
	if err != nil {
		logger.Error(reqCtx, eventClassProduct, "GetMany", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := utils.TransformSlice(products, func(product *models.Product) *models.ProductResponse {
		return product.ToResponse()
	})
	return ctx.JSON(http.StatusOK, response)
}

// GetManyWithPagination handles GET requests to retrieve paginated Products.
func (c *ProductController) GetManyWithPagination(ctx echo.Context) error {
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

	products, pagination, err := c.services.GetManyWithPagination(reqCtx, filter, page, limit)
	if err != nil {
		logger.Error(reqCtx, eventClassProduct, "GetManyWithPagination", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := utils.TransformSlice(products, func(product *models.Product) *models.ProductResponse {
		return product.ToResponse()
	})
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data":       response,
		"pagination": pagination,
	})
}
