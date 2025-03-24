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
    go get github.com/adasarpan404/gobull
    ```

2. Import GoBull in your project

    ```go
    import "github.com/adasarpan404/gobull/queue"
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

## Roadmap 🛣️

✅ Redis-based queue
✅ Delayed jobs
✅ Automatic retries
✅ Multiple workers
⏳ Dead-letter queue
⏳ Priority job scheduling
⏳ Web UI for monitoring jobs

## Contributing 🤝

We welcome contributions! Feel free to:

- Open an issue for bug reports and feature requests.

- Submit a pull request with new features or fixes.

1. Fork the repository

2. Create a feature branch (git checkout -b feature-new)

3. Commit changes (git commit -m "Added new feature")

4. Push to the branch (git push origin feature-new)

5. Open a Pull Request
