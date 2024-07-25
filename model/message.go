package model

type Message struct {
	Content   string `json:"content"`
	Processed bool   `json:"processed"`
}
