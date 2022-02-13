package message

import "encoding/json"

type MessageTypeEvent = uint64

const (
	MTGameMessage MessageTypeEvent = 102
)

type SocketMessage struct {
	MessageType uint64          `json:"msg_type"`
	PayloadType string          `json:"payload_type"`
	Payload     json.RawMessage `json:"payload"`
}
