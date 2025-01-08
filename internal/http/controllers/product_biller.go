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

// ProductBillerController defines the HTTP layer for ProductBiller entities.
type ProductBillerController struct {
	service service.ProductBillerService
}

// NewProductBillerController creates a new instance of ProductBillerController.
func NewProductBillerController(service service.ProductBillerService) *ProductBillerController {
	return &ProductBillerController{
		service: service,
	}
}

// Create handles POST requests to create a new ProductBiller.
func (c *ProductBillerController) Create(ctx echo.Context) error {
	var productBiller *models.CreateProductBillerRequest
	if err := ctx.Bind(productBiller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(productBiller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.service.Create(ctx.Request().Context(), productBiller.ToEntity()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"message": "Product Biller created successfully"})
}

// Update handles PUT requests to update an existing ProductBiller.
func (c *ProductBillerController) Update(ctx echo.Context) error {
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

	if err := c.service.Update(ctx.Request().Context(), uint(id), productBiller.ToEntity()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Product Biller updated successfully"})
}

// Delete handles DELETE requests to remove a ProductBiller.
func (c *ProductBillerController) Delete(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	if err := c.service.Delete(ctx.Request().Context(), uint(id)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Product Biller deleted successfully"})
}

// GetOne handles GET requests to retrieve a single ProductBiller.
func (c *ProductBillerController) GetOne(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID: must be a positive integer")
	}

	productBiller, err := c.service.GetOne(ctx.Request().Context(), uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, productBiller.ToResponse())
}

// GetMany handles GET requests to retrieve multiple ProductBillers based on filters.
func (c *ProductBillerController) GetMany(ctx echo.Context) error {
	filter := make(map[string]interface{})
	if productID := ctx.QueryParam("product_id"); productID != "" {
		filter["product_id"] = productID
	}
	if billerID := ctx.QueryParam("biller_id"); billerID != "" {
		filter["biller_id"] = billerID
	}

	productBillers, err := c.service.GetMany(ctx.Request().Context(), filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := common.TransformSlice(productBillers, func(pb *models.ProductBiller) *models.ProductBillerResponse {
		return pb.ToResponse()
	})
	return ctx.JSON(http.StatusOK, response)
}

// GetManyWithPagination handles GET requests to retrieve paginated ProductBillers.
func (c *ProductBillerController) GetManyWithPagination(ctx echo.Context) error {
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

	productBillers, pagination, err := c.service.GetManyWithPagination(ctx.Request().Context(), filter, page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := common.TransformSlice(productBillers, func(pb *models.ProductBiller) *models.ProductBillerResponse {
		return pb.ToResponse()
	})
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data":       response,
		"pagination": pagination,
	})
}
