package jobs

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Scheduler interface {
	Start()
	Stop()
	ScheduleJob(job job) error
	ScheduleJobWithDelay(job job, delay time.Duration)
	GetQueueLength() int
	IsRunning() bool
}

type scheduler struct {
	jobQueue   chan job
	workers    int
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	retryDelay time.Duration
	maxRetries int
}

func NewScheduler(workers int, queueSize int) Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &scheduler{
		jobQueue:   make(chan job, queueSize),
		workers:    workers,
		ctx:        ctx,
		cancel:     cancel,
		retryDelay: time.Second * 5,
		maxRetries: 3,
	}
}

func (s *scheduler) Start() {
	log.Printf("Starting scheduler with %d workers", s.workers)

	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
}

func (s *scheduler) Stop() {
	log.Println("Stopping scheduler...")
	s.cancel()
	close(s.jobQueue)
	s.wg.Wait()
	log.Println("Scheduler stopped")
}

func (s *scheduler) ScheduleJob(job job) error {
	select {
	case s.jobQueue <- job:
		log.Printf("Job %s scheduled", job.GetID())
		return nil
	case <-s.ctx.Done():
		return fmt.Errorf("scheduler is shutting down")
	default:
		return fmt.Errorf("job queue is full")
	}
}

func (s *scheduler) ScheduleJobWithDelay(job job, delay time.Duration) {
	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()

		select {
		case <-timer.C:
			if err := s.ScheduleJob(job); err != nil {
				log.Printf("Failed to schedule delayed job %s: %v", job.GetID(), err)
			}
		case <-s.ctx.Done():
			return
		}
	}()
}

// worker processes jobs from the queue
func (s *scheduler) worker(id int) {
	defer s.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case job, ok := <-s.jobQueue:
			if !ok {
				log.Printf("Worker %d: job queue closed", id)
				return
			}

			s.executeJob(job, id)

		case <-s.ctx.Done():
			log.Printf("Worker %d: context cancelled", id)
			return
		}
	}
}

// executeJob executes a job with retry logic
func (s *scheduler) executeJob(job job, workerID int) {
	jobCtx, cancel := context.WithTimeout(s.ctx, time.Minute*5)
	defer cancel()

	log.Printf("Worker %d: executing job %s (attempt %d)",
		workerID, job.GetID(), job.GetRetryCount()+1)

	err := job.Execute(jobCtx)

	if err != nil {
		log.Printf("Worker %d: job %s failed: %v", workerID, job.GetID(), err)

		if job.GetRetryCount() < job.GetMaxRetries() {
			job.IncrementRetry()
			log.Printf("Worker %d: retrying job %s in %v (attempt %d/%d)",
				workerID, job.GetID(), s.retryDelay, job.GetRetryCount()+1, job.GetMaxRetries()+1)

			s.ScheduleJobWithDelay(job, s.retryDelay)
		} else {
			log.Printf("Worker %d: job %s exceeded max retries (%d), giving up",
				workerID, job.GetID(), job.GetMaxRetries())
		}
	} else {
		log.Printf("Worker %d: job %s completed successfully", workerID, job.GetID())
	}
}

// GetQueueLength returns the current number of jobs in the queue
func (s *scheduler) GetQueueLength() int {
	return len(s.jobQueue)
}

// IsRunning returns whether the scheduler is currently running
func (s *scheduler) IsRunning() bool {
	select {
	case <-s.ctx.Done():
		return false
	default:
		return true
	}
}
