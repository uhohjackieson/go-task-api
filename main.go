package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	)

// Task represents a single to-do item
type Task struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Done bool `json:"done"`
}

// in-memory storage
var tasks []Task
var mutex = &sync.Mutex{}
var nextID = 1

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			getTasks(w)
		case "POST":
			createTask(w, r)
		default: http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
}

func get Tasks(w http.ResponseWriter) {
	mutex.Lock()
	defer mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
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
	newTask.ID = nextID
	nextID++
	tasks = append(tasks, newTask)
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func main() {
	http.HandleFunc("/tasks", tasksHandler)

	log.Println("Server running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		}
}


