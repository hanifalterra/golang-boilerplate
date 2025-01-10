package v1

import (
	"github.com/labstack/echo/v4"

	"golang-boilerplate/internal/http/controllers"
)

func RegisterProductRoute(e *echo.Group, productController *controllers.ProductController) {
	productGroup := e.Group("/product-billers")
	productGroup.POST("", productController.Create)
	productGroup.PUT("/:id", productController.Update)
	productGroup.DELETE("/:id", productController.Delete)
	productGroup.GET("/:id", productController.GetOne)
	productGroup.GET("/all", productController.GetMany)
	productGroup.GET("", productController.GetManyWithPagination)
}
