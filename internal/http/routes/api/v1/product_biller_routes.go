package v1

import (
	"github.com/labstack/echo/v4"

	"golang-boilerplate/internal/http/controllers"
)

func RegisterProductBillerRoutes(e *echo.Group, productBillerController *controllers.ProductBillerController) {
	productBillerGroup := e.Group("/product-billers")
	productBillerGroup.POST("", productBillerController.Create)
	productBillerGroup.PUT("/:id", productBillerController.Update)
	productBillerGroup.DELETE("/:id", productBillerController.Delete)
	productBillerGroup.GET("/:id", productBillerController.GetOne)
	productBillerGroup.GET("/all", productBillerController.GetMany)
	productBillerGroup.GET("", productBillerController.GetManyWithPagination)
}
