package app

import (
	"template/internal/adapters/inbound/jobs"
	"template/internal/logger"
	"template/packages/cronjob-go"
)

type cronjobs struct {
	logger logger.Logger
}

func newCronJobs(logger logger.Logger) cronjobs {
	return cronjobs{
		logger: logger,
	}
}

func (c *cronjobs) setupCronJobs() {
	c.logger.Info("Setting up cron jobs")

	scheduler := cronjob.NewScheduler(1, 100)
	cronScheduler := cronjob.NewCronScheduler(scheduler)

	cronScheduler.AddJob("0 2 * * *", c.WithLogger(jobs.FakeJob))

	c.logger.Info("Starting cron jobs")
	cronScheduler.Start()
}

func (c *cronjobs) WithLogger(job func(logger logger.Logger)) func() {
	return func() {
		job(c.logger)
	}
}
