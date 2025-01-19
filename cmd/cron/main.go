package main

import (
	"github.com/rs/zerolog/log"

	"golang-boilerplate/internal/cron/config"
	"golang-boilerplate/internal/cron/controllers"
	"golang-boilerplate/internal/cron/usecases"
	"golang-boilerplate/internal/pkg/connections/cacabot"
	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/infrastructure/notification"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/logger"
	"golang-boilerplate/internal/pkg/utils"
)

func main() {
	// Load application configuration
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load application configuration")
	}

	// Initialize logger
	logger := logger.New(
		config.Logger.Level,
		config.App.Name,
		config.App.Version,
		config.Service.Name,
	)

	// Establish database connection
	dbConn, err := db.NewDB(&config.DB, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Database connection failed")
	}

	// Initialize Cacabot client
	cacabotClient := cacabot.NewCacabotClient(
		config.Cacabot.URL,
		config.Cacabot.Username,
		config.Cacabot.Password,
		config.Cacabot.Enabled,
	)

	// Initialize repositories
	productRepo := repositories.NewProductBillerRepository(dbConn)

	// Initialize notification infrastructure
	notif := notification.NewNotification(cacabotClient)

	// Initialize use case layer
	cronUseCase := usecases.NewCronUseCase(productRepo, notif)

	// Initialize controller layer
	cronController := controllers.NewCronController(cronUseCase, logger)

	// Schedule the daily cron job for sending product biller summaries
	cronJob := &utils.CronJob{
		Task: cronController.NotifyProductBillerSummary,
	}
	go cronJob.ScheduleDaily(
		config.Service.NotificationHour,
		config.Service.NotificationMinute,
	)

	// Keep the application running indefinitely
	select {}
}
