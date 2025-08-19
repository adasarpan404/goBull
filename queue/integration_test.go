package queue

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// TestIntegration_ProcessJob connects to a real Redis at REDIS_ADDR (or
// localhost:6379). If Redis isn't available the test is skipped. The test
// starts a Worker, registers a handler that fails once then succeeds, adds a
// job and verifies the handler ran by checking a Redis key set by the handler.
func TestIntegration_ProcessJob(t *testing.T) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available at %s, skipping integration test: %v", redisAddr, err)
	}

	// unique queue name to avoid collisions
	qname := fmt.Sprintf("itest:%d", time.Now().UnixNano())
	q := NewQueue(qname, redisAddr)

	// use the same redis client to inspect keys
	rc := client

	w := NewWorker(q)

	// handler fails first time then succeeds and writes a key when processed
	calls := 0
	handler := func(j *Job) error {
		calls++
		if calls == 1 {
			return fmt.Errorf("simulated failure")
		}
		key := fmt.Sprintf("processed:%s", j.ID)
		// store the attempts observed by the handler
		return rc.Set(context.Background(), key, fmt.Sprintf("attempts:%d", j.Attempts), time.Minute).Err()
	}

	w.RegisterHandler("email", handler)

	// start worker in background with cancellable context
	ctx, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()
	go w.Start(ctx)

	jobID := fmt.Sprintf("integ-%d", time.Now().UnixNano())
	job := Job{ID: jobID, Data: "email", Delay: 0}
	if err := q.AddJob(job); err != nil {
		t.Fatalf("AddJob failed: %v", err)
	}

	// wait up to 10s for the processed key to appear
	found := false
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		val, err := rc.Get(context.Background(), "processed:"+jobID).Result()
		if err == nil {
			t.Logf("found processed key: %s", val)
			found = true
			// stop worker and exit
			cancelWorker()
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	if !found {
		t.Fatalf("job was not processed within timeout")
	}
}
