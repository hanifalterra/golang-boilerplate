package routes

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"

	"golang-boilerplate/internal/http/config"
	"golang-boilerplate/internal/http/controllers"
	v1 "golang-boilerplate/internal/http/routes/api/v1"
	"golang-boilerplate/internal/http/usecases"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/logger"
)

// RegisterRoutes sets up all HTTP routes, middleware, and usecases.
func RegisterRoutes(e *echo.Echo, db *sqlx.DB, log *zerolog.Logger, config *config.Config) {
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
	e.Use(echoprometheus.NewMiddleware(config.Service.Name))
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

	// Initialize UseCases
	productUseCase := usecases.NewProductUseCase(productRepo, uow)
	billerUseCase := usecases.NewBillerUseCase(billerRepo, uow)
	productBillerUseCase := usecases.NewProductBillerUseCase(productBillerRepo, productRepo, billerRepo)

	// Initialize Controllers
	productCtrl := controllers.NewProductController(productUseCase, log)
	billerCtrl := controllers.NewBillerController(billerUseCase, log)
	productBillerCtrl := controllers.NewProductBillerController(productBillerUseCase, log)

	// Register API Version 1 Routes
	apiV1 := e.Group("/api/v1")
	v1.RegisterProductRoute(apiV1, productCtrl)
	v1.RegisterBillerRoute(apiV1, billerCtrl)
	v1.RegisterProductBillerRoute(apiV1, productBillerCtrl)
}
