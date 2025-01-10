package routes

import (
	v1 "golang-boilerplate/internal/http/routes/api/v1"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"golang-boilerplate/internal/http/controllers"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/services"
)

// RegisterRoutes sets up all routes for the application.
func RegisterRoutes(e *echo.Echo, db *sqlx.DB) {
	// Initialize Repositories
	productRepo := repositories.NewProductRepository(db)
	billerRepo := repositories.NewBillerRepository(db)
	productBillerRepo := repositories.NewProductBillerRepository(db)
	uow := repositories.NewUnitOfWork(db)

	// Initialize Services
	productService := services.NewProductService(productRepo, uow)
	billerService := services.NewBillerService(billerRepo, uow)
	productBillerService := services.NewProductBillerService(productBillerRepo, productRepo, billerRepo)

	// Initialize Controllers
	productController := controllers.NewProductController(productService)
	billerController := controllers.NewBillerController(billerService)
	productBillerController := controllers.NewProductBillerController(productBillerService)

	// API group
	api := e.Group("/api")

	// Version 1 routes
	v1Group := api.Group("/v1")
	v1.RegisterProductRoutes(v1Group, productController)
	v1.RegisterBillerRoutes(v1Group, billerController)
	v1.RegisterProductBillerRoutes(v1Group, productBillerController)
}
