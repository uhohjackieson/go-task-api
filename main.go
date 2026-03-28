package main

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

	mutex.Lock()
	tasks = append(tasks, newTask)
	mutex.Unlock()

	json.NewEncoder(w).Encode(newTask)
}

func main() {
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