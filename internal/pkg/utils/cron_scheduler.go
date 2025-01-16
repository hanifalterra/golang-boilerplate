package utils

import (
	"time"
)

type CronJob struct {
	Task func()
}

func (c *CronJob) ScheduleDaily(hour int, minute int) {
	for {
		now := time.Now()
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
		if now.After(nextRun) {
			nextRun = nextRun.Add(24 * time.Hour)
		}

		time.Sleep(time.Until(nextRun))
		c.Task()
	}
}
