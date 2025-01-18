package main

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"golang-boilerplate/internal/pkg/config"
	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/connections/kafka"
	"golang-boilerplate/internal/pkg/connections/redis"
	"golang-boilerplate/internal/pkg/infrastructure/lock"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/logger"
	"golang-boilerplate/internal/worker/controllers"
	"golang-boilerplate/internal/worker/usecases"
)

func main() {
	ctx := context.Background()

	// Load configuration
	appConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize logger
	appLogger := logger.New(appConfig.Logger.Level, appConfig.App.Name, appConfig.App.Version, appConfig.WorkerService.Name)

	// Kafka configuration.
	brokers := "localhost:9092"
	groupID := "transaction-consumer-group"
	topic := "transaction"

	// Initialize Kafka consumer.
	consumer, err := kafka.NewKafkaConsumer(brokers, groupID, topic)
	if err != nil {
		log.Fatal().Msgf("Failed to initialize Kafka consumer: %v", err)
	}

	// Initialize DB connection
	dbConn, err := db.NewDB(&appConfig.DB, appLogger)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	redis, err := redis.NewRedis(context.Background(), &appConfig.Redis)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to redis")
	}

	pbRepo := repositories.NewProductBillerRepository(dbConn)
	lock := lock.NewLock(redis, appConfig.Lock.TTL*time.Millisecond, appConfig.Lock.MaxRetryTime*time.Millisecond, appConfig.Lock.RetryInterval*time.Millisecond)

	// Initialize usecase and controller.
	usecase := usecases.NewTransactionUseCase(pbRepo, lock)
	controller := controllers.NewTransactionController(usecase)

	// Start consuming messages.
	log.Info().Msg("Starting Kafka consumer...")
	consumer.Consume(ctx, controller.HandleMessage)
}
