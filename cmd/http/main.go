package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"golang-boilerplate/internal/http/routes"
	"golang-boilerplate/internal/pkg/config"
	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/logger"
)

func main() {
	// Load configuration
	appConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize logger
	appLogger := logger.New(appConfig.Logger.Level, appConfig.App.Name, appConfig.App.Version, appConfig.HTTPService.Name)

	// Initialize DB connection
	dbConn, err := db.NewDB(&appConfig.DB, appLogger)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	// Initialize Echo server
	server := echo.New()
	routes.RegisterRoutes(server, dbConn, appLogger, appConfig)

	// Run the server in a separate goroutine
	go func() {
		if err := server.Start(":" + appConfig.HTTPService.Port); err != nil {
			appLogger.Fatal().Err(err).Msg("Failed to start the server")
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server, appLogger)
}

// gracefulShutdown handles server shutdown on receiving termination signals
func gracefulShutdown(server *echo.Echo, appLogger *zerolog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	appLogger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error().Err(err).Msg("Error during server shutdown")
	} else {
		appLogger.Info().Msg("Server shutdown completed")
	}
}
