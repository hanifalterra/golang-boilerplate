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

// ProductController defines the HTTP layer for Product entities.
type ProductController struct {
	usecases usecases.ProductUseCase
	logger   *zerolog.Logger
}

// NewProductController creates a new instance of ProductController.
func NewProductController(usecases usecases.ProductUseCase, logger *zerolog.Logger) *ProductController {
	return &ProductController{
		usecases: usecases,
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

	if err := c.usecases.Create(reqCtx, product.ToEntity()); err != nil {
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

	if err := c.usecases.Update(reqCtx, uint(id), product.ToEntity()); err != nil {
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

	if err := c.usecases.Delete(reqCtx, uint(id)); err != nil {
		logger.Error(reqCtx, eventClassProduct, "Delete", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

// FetchOne handles GET requests to retrieve a single Product.
func (c *ProductController) FetchOne(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	product, err := c.usecases.FetchOne(reqCtx, uint(id))
	if err != nil {
		logger.Error(reqCtx, eventClassProduct, "FetchOne", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, product.ToResponse())
}

// FetchMany handles GET requests to retrieve multiple Products based on filters.
func (c *ProductController) FetchMany(ctx echo.Context) error {
	reqCtx, logger := logger.NewAppLoggerEcho(ctx, c.logger)

	filter := make(map[string]interface{})
	if label := ctx.QueryParam("label"); label != "" {
		filter["label"] = label
	}

	products, err := c.usecases.FetchMany(reqCtx, filter)
	if err != nil {
		logger.Error(reqCtx, eventClassProduct, "FetchMany", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := utils.TransformSlice(products, func(product *models.Product) *models.ProductResponse {
		return product.ToResponse()
	})
	return ctx.JSON(http.StatusOK, response)
}

// FetchManyWithPagination handles GET requests to retrieve paginated Products.
func (c *ProductController) FetchManyWithPagination(ctx echo.Context) error {
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

	products, pagination, err := c.usecases.FetchManyWithPagination(reqCtx, filter, page, limit)
	if err != nil {
		logger.Error(reqCtx, eventClassProduct, "FetchManyWithPagination", err.Error())
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
