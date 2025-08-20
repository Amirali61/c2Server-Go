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
func recvAcknowledge(r *http.Response) {
	ack := make([]byte, 2)
	r.Body.Read(ack)
	fmt.Println(string(ack))

}
func main() {
	for {
		resp := sendBeacon("1234")
		recvAcknowledge(resp)
		time.Sleep(4 * time.Second)
	}
}
