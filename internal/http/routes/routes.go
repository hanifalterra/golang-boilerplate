package routes

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"

	"golang-boilerplate/internal/http/controllers"
	v1 "golang-boilerplate/internal/http/routes/api/v1"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/logger"
	"golang-boilerplate/internal/pkg/services"
)

const serviceName = "golang-boilerplate-http"

// RegisterRoutes sets up all HTTP routes, middleware, and services.
func RegisterRoutes(e *echo.Echo, db *sqlx.DB, log *zerolog.Logger) {
	// General Middleware Configuration
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))
	e.Use(middleware.RequestID())
	e.Use(echoprometheus.NewMiddleware(serviceName))
	e.GET("/metrics", echoprometheus.NewHandler())

	// Request Logging Middleware
	requestLogger := logger.NewLoggerMiddleware(log)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		HandleError:   true, // Log errors from the response
		LogError:      true,
		LogRequestID:  true,
		LogLatency:    true,
		LogValuesFunc: requestLogger.LogRequest, // Custom log function
	}))

	// Initialize Unit of Work
	uow := repositories.NewUnitOfWork(db)

	// Initialize Repositories
	productRepo := repositories.NewProductRepository(db)
	billerRepo := repositories.NewBillerRepository(db)
	productBillerRepo := repositories.NewProductBillerRepository(db)

	// Initialize Services
	productService := services.NewProductService(productRepo, uow)
	billerService := services.NewBillerService(billerRepo, uow)
	productBillerService := services.NewProductBillerService(productBillerRepo, productRepo, billerRepo)

	// Initialize Controllers
	productCtrl := controllers.NewProductController(productService, log)
	billerCtrl := controllers.NewBillerController(billerService, log)
	productBillerCtrl := controllers.NewProductBillerController(productBillerService, log)

	// Register API Version 1 Routes
	apiV1 := e.Group("/api/v1")
	v1.RegisterProductRoute(apiV1, productCtrl)
	v1.RegisterBillerRoute(apiV1, billerCtrl)
	v1.RegisterProductBillerRoute(apiV1, productBillerCtrl)
}
