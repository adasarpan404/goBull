package queue

import (
	"context"
	"fmt"
	"time"
)

type Worker struct {
	queue    *Queue
	handlers map[string]func(*Job) error
}

func NewWorker(queue *Queue) *Worker {
	return &Worker{
		queue:    queue,
		handlers: make(map[string]func(*Job) error),
	}
}

func (w *Worker) RegisterHandler(jobType string, handler func(*Job) error) {
	w.handlers[jobType] = handler
}

// Start runs the worker loop until the provided context is cancelled.
func (w *Worker) Start(ctx context.Context) {
	fmt.Println("Worker Started......")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker stopped by context")
			return
		default:
		}

		job, err := w.queue.GetJob(ctx)

		if err != nil {
			fmt.Println("Error fetching job:", err)
			continue
		}

		handler, exists := w.handlers[job.Data]
		if !exists {
			fmt.Println("No handler found for job:", job.Data)
			continue
		}
		fmt.Println("Processing job:", job.ID)
		err = handler(job)
		if err != nil {
			fmt.Println("Error processing job:", err)
			// Retry logic
			job.Attempts++
			if job.Attempts < 3 {
				fmt.Println("Retrying job:", job.ID)
				time.Sleep(2 * time.Second)
				_ = w.queue.AddJob(*job)
			}
		}
	}
}
