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

1.Install GoBull as a Module

```bash
go get github.com/yourusername/gobull
```

2.Import GoBull in your project

```go
import "github.com/yourusername/gobull/queue"
```

