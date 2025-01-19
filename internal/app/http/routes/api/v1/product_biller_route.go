package v1

import (
	"github.com/labstack/echo/v4"

	"golang-boilerplate/internal/app/http/controllers"
)

func RegisterProductBillerRoute(e *echo.Group, productBillerController *controllers.ProductBillerController) {
	productBillerGroup := e.Group("/product-billers")
	productBillerGroup.POST("", productBillerController.Create)
	productBillerGroup.PUT("/:id", productBillerController.Update)
	productBillerGroup.DELETE("/:id", productBillerController.Delete)
	productBillerGroup.GET("/:id", productBillerController.FetchOne)
	productBillerGroup.GET("/all", productBillerController.FetchMany)
	productBillerGroup.GET("", productBillerController.FetchManyWithPagination)
}
