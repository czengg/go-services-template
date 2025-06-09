package jobs

import (
	"fmt"
)

func SetupCronJobs() {
	fmt.Println("Starting cron jobs")

	scheduler := NewScheduler(1, 100)
	cronScheduler := NewCronScheduler(scheduler)

	cronScheduler.AddJob("0 2 * * *", func() {
		fmt.Println("job running")
	})

}
