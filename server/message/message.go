package message

import "encoding/json"

type MessageTypeEvent = uint64

const (
	MTNewEvents MessageTypeEvent = 100
	MTEventSync MessageTypeEvent = 101
)

type EventSync struct {
	GameTick uint64          `json:"gametick"`
	Events   json.RawMessage `json:"events"`
}

type SocketMessage struct {
	MessageType uint64          `json:"msg_type"`
	Payload     json.RawMessage `json:"payload"`
}
