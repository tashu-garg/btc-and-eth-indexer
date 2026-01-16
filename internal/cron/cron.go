package cron

import (
	"time"
)

// Job defines a function to be executed
type Job func()

// Scheduler manages periodic jobs
type Scheduler struct {
	stopCh chan struct{}
}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		stopCh: make(chan struct{}),
	}
}

// AddJob starts a job running at the specified interval
func (s *Scheduler) AddJob(interval time.Duration, job Job) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				job()
			case <-s.stopCh:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops all jobs
func (s *Scheduler) Stop() {
	close(s.stopCh)
}
