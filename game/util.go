package game

import (
	"encoding/json"

	"github.com/emyrk/grow/game/events"
)

func NewEventMsgPayload(evts []events.Event) (GameMessageType, []byte) {
	msg := NewEvents{
		Eventlist: evts,
	}
	data, _ := json.Marshal(msg)
	return msg.Type(), data
}
