package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hibiken/asynq"
)

type InputParameter struct {
	Operator1 int64
	Operator2 int64
}

func main() {
	client := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr: "localhost:6379",
		},
	)

	http.HandleFunc("/schedule-task-high-priority", func(w http.ResponseWriter, r *http.Request) {
		arguments, err := json.Marshal(InputParameter{Operator1: 2, Operator2: 3})
		if err != nil {
			log.Fatal(err)
		}

		task := asynq.NewTask("addition-task", arguments)
		info, err := client.Enqueue(task, asynq.Queue("high-priority"))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[*] Successfully enqueued task: %+v", info)
	})

	http.HandleFunc("/schedule-task-low-priority", func(w http.ResponseWriter, r *http.Request) {
	})

	http.HandleFunc("/schedule-task", func(w http.ResponseWriter, r *http.Request) {
	})

	log.Printf("Starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
