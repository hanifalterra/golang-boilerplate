package v1

import (
	"github.com/labstack/echo/v4"

	"golang-boilerplate/internal/app/http/controllers"
)

func RegisterBillerRoute(e *echo.Group, billerController *controllers.BillerController) {
	billerGroup := e.Group("/product-billers")
	billerGroup.POST("", billerController.Create)
	billerGroup.PUT("/:id", billerController.Update)
	billerGroup.DELETE("/:id", billerController.Delete)
	billerGroup.GET("/:id", billerController.FetchOne)
	billerGroup.GET("/all", billerController.FetchMany)
	billerGroup.GET("", billerController.FetchManyWithPagination)
}
