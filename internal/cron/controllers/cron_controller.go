package controllers

import (
	"fmt"
	"golang-boilerplate/internal/cron/services"
)

type CronController struct {
	Service *services.ProductBillerService
}

func NewCronController(service *services.ProductBillerService) *CronController {
	return &CronController{Service: service}
}

func (c *CronController) RunDailyTask() {
	if err := c.Service.ProcessInactiveBillers(); err != nil {
		// Log the error, but ensure it does not stop execution
		fmt.Printf("Error processing inactive billers: %v\n", err)
	}
}
