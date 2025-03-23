package queue

import "fmt"

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

func (w *Worker) start() {
	fmt.Println("Worker Started......")

	for {
		job, err := w.queue.GetJob()

		if err != nil {
			fmt.Println()
		}
	}
}
