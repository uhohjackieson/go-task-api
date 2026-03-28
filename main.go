package main

var taskChannel = make(chan Task, 100)

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	tasks []Task
	mutex sync.Mutex
)

func startWorker() {
	go func() {
		for t := range taskChannel {
			mutex.Lock()
			tasks = append(tasks, t)
			mutex.Unlock()
			log.Println("Processed task asynchronously:", t.Name)
			}
		}()
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Send task to the channel to be processed asynchronously
	taskChannel <- newTask

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Task queued for processing",
		"name":    newTask.Name,
	})
}

func main() {
	startWorker()
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTasks(w, r)
		} else if r.Method == http.MethodPost {
			createTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
