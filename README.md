# GoBull - A Redis-backed Message Queue in Go

GoBull is a lightweight, high-performance message queue built in Go, inspired by BullMQ. It provides job **scheduling**, **delayed execution**, **automatic retries**, and **distributed processing** using Redis.

## Features 🚀

- **Job Queuing** - Push jobs to a queue and process them asynchronously.

- **Delayed Jobs** - Schedule jobs to run after a specific delay.

- **Automatic Retries** - Failed jobs are retried with exponential backoff.

- **Worker System** - Register workers to handle specific job types.

- **Distributed Processing** - Scale workers across multiple instances.

- **Persistent Storage** - Uses Redis to store and manage jobs.

## Integrating GoBull in an Existing Project 🛠️

To use GoBull in your existing Go project, follow these steps:

1. Install GoBull as a Module

    ```bash
    go get github.com/yourusername/gobull
    ```

2. Import GoBull in your project

    ```go
    import "github.com/yourusername/gobull/queue"
    ```

3. Create a Queue in Your Application

    ```go
    q:= queue.NewQueue("my_tasks", "localhost:6379")
    ```

4. Add Jobs from Your Existing Code

    ```go
    job := queue.Job{ID: "123", Data: "send_email", Delay: 10}
    q.AddJob(job)
    ```

5. Start a Worker to Process Jobs

    ```go
    worker := queue.NewWorker(q)
    worker.RegisterHandler("send_email", func(job *queue.Job) error {
        fmt.Println("Processing email job:", job.ID)
        return nil
    })
    worker.Start()
    ```

Now, your existing project can use GoBull for background job processing! 🎉

## API Reference 📌

### Queue Methods

| Method               | Description                             |
|----------------------|-----------------------------------------|
| `AddJob(job Job) error` | Adds a job to the queue                 |
| `GetJob() (*Job, error)` | Retrieves and removes a job from the queue |

## Worker Methods

| Method                                               | Description                             |
|------------------------------------------------------|-----------------------------------------|
| `RegisterHandler(jobType string, handler func(*Job) error)` | Registers a worker for a job type       |
| `Start()`                                            | Starts consuming jobs                   |
