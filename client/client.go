package main

import (
	"bytes"
	"c2-server/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func sendBeacon(id string) *http.Response {
	var beacon models.NewBeacon
	beacon.ID = id

	jsonData, err := encodeBeacon(beacon)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post("http://localhost:5000/beacon", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	return resp
}
func encodeBeacon(beacon models.NewBeacon) ([]byte, error) {
	return json.Marshal(beacon)
}
func recvTask(r *http.Response) {
	var task models.NewTask
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		fmt.Println("Error receiving task from server: " + err.Error())
		return
	}
	fmt.Println("[server] ---> Task: " + string(task.Command))
}
func main() {
	for {
		resp := sendBeacon("1234")
		recvTask(resp)
		time.Sleep(4 * time.Second)
	}
}
