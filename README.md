# GoBull

GoBull is a small Redis-backed job queue written in Go. It aims to be
easy to use in development and production, with support for delayed jobs,
automatic retries, simple worker registration, and a small test surface.

This repository contains the core `queue` package and a couple of
examples to demonstrate usage.

## Quick overview

- Queue: create with `queue.NewQueue(name, redisAddr)`; add jobs with
    `AddJob(Job)` and fetch with `GetJob()`.
- Worker: `NewWorker(queue)` → `RegisterHandler(jobType, handler)` →
    `Start()` to begin processing.

The API is intentionally small so it is easy to reason about and test.

## Getting started (local)

Prerequisites

- Go (module-aware) matching the project `go.mod` (this repo uses Go 1.23+)
- Redis for running examples; tests use an in-memory Redis (miniredis)

Clone and prepare:

```bash
git clone https://github.com/adasarpan404/goBull.git
cd goBull
go mod tidy
```

Run tests (this uses miniredis so you don't need a running Redis):

```bash
go test ./... -v
```

## Usage examples

Library usage (example):

```go
package main

import (
        "fmt"
        "time"

        "github.com/adasarpan404/goBull/queue"
)

func main() {
        q := queue.NewQueue("tasks", "localhost:6379")

        // add a job
        job := queue.Job{ID: "1", Data: "email", Delay: 0}
        if err := q.AddJob(job); err != nil {
                panic(err)
        }

        // start a worker
        w := queue.NewWorker(q)
        w.RegisterHandler("email", func(j *queue.Job) error {
                fmt.Println("processing job", j.ID)
                time.Sleep(1 * time.Second)
                return nil
        })

        // run worker (blocking)
        w.Start()
}
```

Notes about the `examples/` package

- The files in `examples/` provide helper functions (e.g. `RunProducer`,
    `RunConsumer`) rather than `main` functions so they build during tests.
    To run them directly, call those functions from a small `main` program
    like the snippet above.

## Tests and development

- Unit tests use `github.com/alicebob/miniredis/v2` so they run without
    a real Redis server.
- To run a single package tests:

```bash
go test ./queue -v
```

If you want to run examples against a real Redis instance, start Redis
locally (e.g. via Docker) and run your small `main` program.

## Next improvements (roadmap)

- Dead-letter queue (DLQ) and configurable retry/backoff behavior
- Prometheus metrics and observability
- A tiny web UI for inspecting queues and DLQ
- Official Docker image and docker-compose for local development

## Contributing

Contributions are welcome. Please open issues for feature requests or
bugs. If you send a pull request, include tests and a short description
of the change.

## Where to go from here

- Run the tests: `go test ./... -v`
- Try the snippet above with a local Redis
- I can open PRs to add DLQ support, metrics, CI, or a CLI if you want—tell me which to prioritize.
