package controllers

import (
	"context"
	"fmt"

	"golang-boilerplate/internal/cron/usecases"
)

type CronController struct {
	usecase *usecases.ProductBillerUseCase
}

func NewCronController(usecase *usecases.ProductBillerUseCase) *CronController {
	return &CronController{usecase: usecase}
}

func (c *CronController) RunDailyTask() {
	if err := c.usecase.ProcessCountProductBillers(context.Background()); err != nil {
		// Log the error, but ensure it does not stop execution
		fmt.Printf("Error processing inactive billers: %v\n", err)
	}
}
