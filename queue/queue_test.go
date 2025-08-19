package queue

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// TestAddAndGetJob ensures AddJob stores a job on the list and GetJob retrieves it.
func TestAddAndGetJob(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	q := NewQueue("tasks", mr.Addr())

	job := Job{ID: "1", Data: "email", Delay: 0}
	if err := q.AddJob(job); err != nil {
		t.Fatalf("AddJob failed: %v", err)
	}

	got, err := q.GetJob()
	if err != nil {
		t.Fatalf("GetJob failed: %v", err)
	}

	if got.ID != job.ID {
		t.Fatalf("expected job ID %s, got %s", job.ID, got.ID)
	}
	if got.Data != job.Data {
		t.Fatalf("expected job Data %s, got %s", job.Data, got.Data)
	}
}

// TestAddJobDelayed ensures jobs with Delay are added to the delayed zset.
func TestAddJobDelayed(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	q := NewQueue("tasks", mr.Addr())

	job := Job{ID: "2", Data: "report", Delay: 5}
	if err := q.AddJob(job); err != nil {
		t.Fatalf("AddJob (delayed) failed: %v", err)
	}

	// Use a direct redis client to inspect the zset created by AddJob
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	zcard, err := client.ZCard(ctx, "tasks:delayed").Result()
	if err != nil {
		t.Fatalf("ZCard failed: %v", err)
	}
	if zcard != 1 {
		t.Fatalf("expected 1 delayed job in zset, got %d", zcard)
	}
}

// TestWorkerRegisterHandler verifies that RegisterHandler stores the handler.
func TestWorkerRegisterHandler(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	q := NewQueue("tasks", mr.Addr())
	w := NewWorker(q)

	handler := func(j *Job) error { return nil }
	w.RegisterHandler("email", handler)

	if _, ok := w.handlers["email"]; !ok {
		t.Fatalf("handler for 'email' not registered")
	}
}
