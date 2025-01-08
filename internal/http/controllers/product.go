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

// ProductController defines the HTTP layer for Product entities.
type ProductController struct {
	service service.ProductService
}

// NewProductController creates a new instance of ProductController.
func NewProductController(service service.ProductService) *ProductController {
	return &ProductController{
		service: service,
	}
}

// Create handles POST requests to create a new Product.
func (c *ProductController) Create(ctx echo.Context) error {
	var product *models.CreateProductRequest
	if err := ctx.Bind(product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.service.Create(ctx.Request().Context(), product.ToEntity()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"message": "Product created successfully"})
}

// Update handles PUT requests to update an existing Product.
func (c *ProductController) Update(ctx echo.Context) error {
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

	if err := c.service.Update(ctx.Request().Context(), uint(id), product.ToEntity()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Product updated successfully"})
}

// Delete handles DELETE requests to remove a Product.
func (c *ProductController) Delete(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	if err := c.service.Delete(ctx.Request().Context(), uint(id)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

// GetOne handles GET requests to retrieve a single Product.
func (c *ProductController) GetOne(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	product, err := c.service.GetOne(ctx.Request().Context(), uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, product.ToResponse())
}

// GetMany handles GET requests to retrieve multiple Products based on filters.
func (c *ProductController) GetMany(ctx echo.Context) error {
	filter := make(map[string]interface{})
	if label := ctx.QueryParam("label"); label != "" {
		filter["label"] = label
	}

	products, err := c.service.GetMany(ctx.Request().Context(), filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := common.TransformSlice(products, func(product *models.Product) *models.ProductResponse {
		return product.ToResponse()
	})
	return ctx.JSON(http.StatusOK, response)
}

// GetManyWithPagination handles GET requests to retrieve paginated Products.
func (c *ProductController) GetManyWithPagination(ctx echo.Context) error {
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

	products, pagination, err := c.service.GetManyWithPagination(ctx.Request().Context(), filter, page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := common.TransformSlice(products, func(product *models.Product) *models.ProductResponse {
		return product.ToResponse()
	})
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data":       response,
		"pagination": pagination,
	})
}
