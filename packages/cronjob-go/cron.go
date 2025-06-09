package cronjob

import (
	"time"

	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	scheduler Scheduler
	cron      *cron.Cron
}

func NewCronScheduler(scheduler Scheduler) *CronScheduler {
	return &CronScheduler{
		scheduler: scheduler,
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithLocation(time.UTC),
		),
	}
}

func (cs *CronScheduler) Start() {
	cs.cron.Start()
}

func (cs *CronScheduler) Stop() {
	cs.cron.Stop()
}

func (cs *CronScheduler) AddJob(spec string, jobFunc func()) {
	cs.cron.AddFunc(spec, jobFunc)
}
