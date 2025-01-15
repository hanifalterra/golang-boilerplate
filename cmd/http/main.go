package main

import (
	"github.com/labstack/echo/v4"

	"golang-boilerplate/internal/http/routes"
	"golang-boilerplate/internal/pkg/config"
	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/logger"
)

const EVENTCLASS = "service"
const SERVICENAME = "admin"

func main() {
	config, _ := config.NewConfig()
	log := logger.New(config, SERVICENAME)
	// Initialize DB connection
	db, _ := db.NewDB(&config.DB, log)

	// Initialize Echo
	e := echo.New()

	// Register Routes
	routes.RegisterRoutes(e, db, log)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
