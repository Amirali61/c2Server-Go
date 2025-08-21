package models

type NewBeacon struct {
	ID string `json:"id"`
}

type NewTask struct {
	ID      string `json:"id"`
	Command string `json:"command"`
}
