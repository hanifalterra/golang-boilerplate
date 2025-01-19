package controllers

import (
	"context"

	"github.com/rs/zerolog"

	"golang-boilerplate/internal/cron/usecases"
	"golang-boilerplate/internal/pkg/logger"
)

type CronController struct {
	usecase *usecases.CronUseCase
	logger  *zerolog.Logger
}

func NewCronController(usecase *usecases.CronUseCase, logger *zerolog.Logger) *CronController {
	return &CronController{
		usecase: usecase,
		logger:  logger,
	}
}

const eventClassCron = "controller.cron"

func (c *CronController) NotifyProductBillerSummary() {
	// Create a logger with a contextualized application logger
	ctx, logger := logger.NewAppLogger(context.Background(), c.logger)

	// Attempt to execute the use case and handle any errors
	if err := c.usecase.NotifyProductBillerSummary(ctx); err != nil {
		logger.Error(ctx, eventClassCron, "NotifyProductBillerSummary", err.Error())
	}
}
