package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Job struct {
	ID        string    `json:"id"`
	TimeStamp time.Time `json:"timestamp"`
	Data      string    `json:"data"`
	Delay     int       `json:"delay"`
	Attempts  int       `json:"attempts"`
}

type Queue struct {
	client *redis.Client
	name   string
	ctx    context.Context
}

func NewQueue(name string, redisAddr string) *Queue {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &Queue{
		client: rdb,
		name:   name,
		ctx:    context.Background(),
	}
}

func (q *Queue) AddJob(job Job) error {
	job.TimeStamp = time.Now()

	jobJSON, err := json.Marshal(job)

	if err != nil {
		return err
	}

	if job.Delay > 0 {
		return q.client.ZAdd(q.ctx, q.name+":delayed", redis.Z{
			Score:  float64(time.Now().Unix() + int64(job.Delay)),
			Member: jobJSON,
		}).Err()
	}

	return q.client.LPush(q.ctx, q.name, jobJSON).Err()
}

func (q *Queue) GetJob(ctx context.Context) (*Job, error) {
	jobstr, err := q.client.BRPop(ctx, 0*time.Second, q.name).Result()
	if err != nil {
		return nil, err
	}
	var job Job
	err = json.Unmarshal([]byte(jobstr[1]), &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
