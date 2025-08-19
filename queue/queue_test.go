package queue

import (
	"context"
	"fmt"
	"sync"
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

	got, err := q.GetJob(context.Background())
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

// TestRetryLogic runs a worker whose handler fails the first time and succeeds the second.
// It verifies the job is retried (Attempts increments) and eventually processed.
func TestRetryLogic(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	q := NewQueue("tasks", mr.Addr())
	w := NewWorker(q)

	var mu sync.Mutex
	calls := 0
	done := make(chan int, 1)

	handler := func(j *Job) error {
		mu.Lock()
		calls++
		callNum := calls
		mu.Unlock()

		if callNum == 1 {
			// simulate failure on first attempt
			return fmt.Errorf("simulated failure")
		}

		// on success, send the Attempts value observed by the handler
		done <- j.Attempts
		return nil
	}

	w.RegisterHandler("retry", handler)

	// run worker in background with cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go w.Start(ctx)

	// add job
	job := Job{ID: "r1", Data: "retry", Delay: 0}
	if err := q.AddJob(job); err != nil {
		t.Fatalf("AddJob failed: %v", err)
	}

	select {
	case attempts := <-done:
		if attempts < 1 {
			t.Fatalf("expected attempts>=1 after retry, got %d", attempts)
		}
		// stop worker
		cancel()
	case <-time.After(6 * time.Second):
		cancel()
		t.Fatalf("timeout waiting for job to be retried and processed")
	}
}

// TestDelayedJobScore asserts that delayed jobs are stored in the delayed zset
// with a score approximately equal to now + Delay (in seconds).
func TestDelayedJobScore(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	q := NewQueue("tasks", mr.Addr())
	delay := 3
	job := Job{ID: "d1", Data: "delayed", Delay: delay}
	before := time.Now().Unix()
	if err := q.AddJob(job); err != nil {
		t.Fatalf("AddJob (delayed) failed: %v", err)
	}

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// fetch zset members with scores
	res, err := client.ZRangeWithScores(ctx, "tasks:delayed", 0, -1).Result()
	if err != nil {
		t.Fatalf("ZRangeWithScores failed: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 delayed job in zset, got %d", len(res))
	}

	score := int64(res[0].Score)
	expected := before + int64(delay)
	// allow a small delta because time advanced between calls
	if score < expected || score > expected+2 {
		t.Fatalf("expected score around %d, got %d", expected, score)
	}
}
