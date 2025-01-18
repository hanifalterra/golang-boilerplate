package main

import (
	"github.com/rs/zerolog/log"

	"golang-boilerplate/internal/cron/controllers"
	"golang-boilerplate/internal/cron/usecases"
	"golang-boilerplate/internal/pkg/config"
	"golang-boilerplate/internal/pkg/connections/cacabot"
	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/infrastructure/notification"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/logger"
	"golang-boilerplate/internal/pkg/utils"
)

func main() {
	// Load configuration
	appConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize logger
	appLogger := logger.New(appConfig.Logger.Level, appConfig.App.Name, appConfig.App.Version, appConfig.CronService.Name)

	// Initialize DB connection
	dbConn, err := db.NewDB(&appConfig.DB, appLogger)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	cacabotClient := cacabot.NewCacabotClient(appConfig.Cacabot.URL, appConfig.Cacabot.Username, appConfig.Cacabot.Password, appConfig.Cacabot.Enabled)

	productRepo := repositories.NewProductBillerRepository(dbConn)
	notif := notification.NewNotification(cacabotClient)
	productUseCase := usecases.NewProductBillerUseCase(productRepo, notif)
	cronCtrl := controllers.NewCronController(productUseCase)

	cron := &utils.CronJob{Task: cronCtrl.RunDailyTask}
	go cron.ScheduleDaily(9, 0)

	// Keep the app running
	select {}
}
