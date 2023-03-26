package service

import "time"

type Sender struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type Location struct {
	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`
}
type Message struct {
	Id         string `json:"id"`
	ClientId   string `json:"clientId"`
	Content    string `json:"content"`
	Sender     `json:"sender"`
	SentAt     time.Time `json:"sentAt"`
	ReceivedAt time.Time `json:"receivedAt"`
	Location   `json:"location"`
}
type MessagePayload struct {
	Type    string `json:"type"`
	Message `json:"message"`
}
type KafkaMessage struct {
	Id      string         `json:"id"`
	Payload MessagePayload `json:"payload"`
}
