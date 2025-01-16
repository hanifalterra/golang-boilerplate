package main

import (
	controller "golang-boilerplate/internal/cron/controllers"
	"golang-boilerplate/internal/cron/services"
	"golang-boilerplate/internal/pkg/config"
	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/infrastructure/notification"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/logger"
	"golang-boilerplate/internal/pkg/utils"

	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration
	appConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize logger
	appLogger := logger.New(appConfig, appConfig.HTTPService.Name)

	// Initialize DB connection
	dbConn, err := db.NewDB(&appConfig.DB, appLogger)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	productRepo := repositories.NewProductBillerRepository(dbConn)
	telegramNotifier := notification.NewTelegramNotifier("https://api.telegram.org/bot<TOKEN>/sendMessage")
	productService := services.NewProductBillerService(productRepo, telegramNotifier)
	cronCtrl := controller.NewCronController(productService)

	cron := &utils.CronJob{Task: cronCtrl.RunDailyTask}
	go cron.ScheduleDaily(9, 0)

	// Keep the app running
	select {}
}
