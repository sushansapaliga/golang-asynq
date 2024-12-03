package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"time"

	"github.com/hibiken/asynq"
)

type InputParameter struct {
	Operator1 int64
	Operator2 int64
}

// Define a type for the context key to avoid key collisions
type contextKey string

const workerHeaderKey contextKey = "x-worker-name"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{
			BaseContext: SetWorkerBaseContext,
			Queues: map[string]int{
				"high-priority":   7,
				"normal-priority": 2,
				"low-priority":    1,
			},
			Concurrency: 1,
		},
	)

	mux := asynq.NewServeMux()
	mux.Use(MiddlewareLogTaskEntryAndExit)
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

func SetWorkerBaseContext() context.Context {
	ctx := context.Background()
	randLength := 10
	charSet := "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
	workerName := "regular-worker-" + StringWithCharset(randLength, charSet)
	return context.WithValue(ctx, workerHeaderKey, workerName)
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset)-1)]
	}
	return string(b)
}

// GetWorkerHeader retrieves the worker name from the context
func GetWorkerHeader(ctx context.Context) string {
	if workerName, ok := ctx.Value(workerHeaderKey).(string); ok {
		return workerName
	}
	return ""
}

func MiddlewareLogTaskEntryAndExit(next asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		slog.Info(fmt.Sprintf("Starting task `%s` by worker: `%s`", t.Type(), GetWorkerHeader(ctx)))
		err := next.ProcessTask(ctx, t)
		slog.Info(fmt.Sprintf("Finished task `%s` by worker: `%s`", t.Type(), GetWorkerHeader(ctx)))
		return err
	})
}
