package main

import (
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

func main() {
	q := queue.NewQueue("tasks", "localhost:6379")
	worker := queue.NewWorker(q)

	worker.RegisterHandler("email", sendEmail)
	worker.RegisterHandler("report", generateReport)

	worker.Start()
}
