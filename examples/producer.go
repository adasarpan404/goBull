package examples

import (
	"fmt"

	"github.com/adasarpan404/goBull/queue"
)

// RunProducer demonstrates how to add jobs to the queue. It's not a main
// function so the package can be built during tests.
func RunProducer() {
	q := queue.NewQueue("tasks", "localhost:6379")

	// job1 := queue.Job{ID: "1", Data: "email", Delay: 0}
	job2 := queue.Job{ID: "2", Data: "report", Delay: 1}

	_ = q.AddJob(job2)
	fmt.Println("Jobs added")
}
