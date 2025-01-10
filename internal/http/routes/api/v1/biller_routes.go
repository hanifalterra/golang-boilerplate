package v1

import (
	"github.com/labstack/echo/v4"

	"golang-boilerplate/internal/http/controllers"
)

func RegisterBillerRoutes(e *echo.Group, billerController *controllers.BillerController) {
	billerGroup := e.Group("/product-billers")
	billerGroup.POST("", billerController.Create)
	billerGroup.PUT("/:id", billerController.Update)
	billerGroup.DELETE("/:id", billerController.Delete)
	billerGroup.GET("/:id", billerController.GetOne)
	billerGroup.GET("/all", billerController.GetMany)
	billerGroup.GET("", billerController.GetManyWithPagination)
}
