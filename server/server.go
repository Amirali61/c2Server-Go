package main

import (
	"c2-server/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var Tasks []models.NewTask

func recvBeacon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := decodeBeacon(r)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("Received ID:", req.ID)
	w.Write([]byte("OK"))
}

func decodeBeacon(r *http.Request) (models.NewBeacon, error) {
	var req models.NewBeacon
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return models.NewBeacon{}, err
	}
	return req, nil
}

func runServer() {
	http.HandleFunc("/beacon", recvBeacon)
	fmt.Println("Server started on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func main() {
	runServer()
}
