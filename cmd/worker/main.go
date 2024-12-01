package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type InputParameter struct {
	Operator1 int64
	Operator2 int64
}

func main() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{
			Queues: map[string]int{
				"high-priority":   7,
				"normal-priority": 2,
				"low-priority":    1,
			},
			Concurrency: 1,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("addition-task", additionTask)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}

func additionTask(ctx context.Context, t *asynq.Task) error {
	var arguments InputParameter

	if err := json.Unmarshal(t.Payload(), &arguments); err != nil {
		return err
	}
	time.Sleep(30 * time.Second)
	log.Printf(" [*] Addition Result is %d", arguments.Operator1+arguments.Operator2)
	return nil
}
