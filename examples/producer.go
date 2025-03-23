package main

import (
	"fmt"

	"github.com/adasarpan404/goBull/queue"
)

func main() {
	q := queue.NewQueue("tasks", "localhost:6379")

	// job1 := queue.Job{ID: "1", Data: "email", Delay: 0}
	job2 := queue.Job{ID: "2", Data: "report", Delay: 1}

	// err := q.AddJob(job1)
	// if err != nil {
	// 	fmt.Println("Failed to add job:", err)
	// }
	err := q.AddJob(job2)
	if err != nil {
		fmt.Println("Failed to add job:", err)
	}
	fmt.Println("Jobs added")
}
