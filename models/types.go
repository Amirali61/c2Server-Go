package models

import "time"

type NewBeacon struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
}

type NewTask struct {
	ID      string `json:"id"`
	Command string `json:"command"`
}

type NewAnswer struct {
	ID     string `json:"id"`
	Answer string `json:"answer"`
}

type Agent struct {
	ID       string
	Hostname string
	LastSeen time.Time
	Tasks    []string
}
