package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/adasarpan404/goBull/queue"
)

func sendEmail(job *queue.Job) error {
	fmt.Println("Sending email:", job.ID)
	time.Sleep(2 * time.Second)
	fmt.Println("Email sent!")
	return nil
}

func generateReport(job *queue.Job) error {
	fmt.Println("Generating report:", job.ID)
	time.Sleep(3 * time.Second)
	fmt.Println("Report generated!")
	return nil
}

// RunConsumer demonstrates starting a worker. It's not a main function so the
// package builds during `go test`.
func RunConsumer() {
	q := queue.NewQueue("tasks", "localhost:6379")
	worker := queue.NewWorker(q)

	worker.RegisterHandler("email", sendEmail)
	worker.RegisterHandler("report", generateReport)

	// run with a cancellable context for demo purposes
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.Start(ctx)

	// demo run for a short time
	time.Sleep(5 * time.Second)
	cancel()
}
