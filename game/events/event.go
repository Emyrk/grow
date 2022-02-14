package events

import (
	"encoding/json"
	world2 "github.com/emyrk/grow/game/world"

	"golang.org/x/xerrors"

	"github.com/rs/zerolog"
)

type EventType string

const (
	LeftClickEvent = "left-click"
	PlayerJoined   = "player-join"
)

type marshalStruct struct {
	EventType EventType       `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

type EventList []Event

func (l EventList) MarshalJSON() ([]byte, error) {
	return MarshalJsonEvents([]Event(l))
}

func (l *EventList) UnmarshalJSON(data []byte) error {
	evts, err := UnmarshalJsonEvents(data)
	if err != nil {
		return err
	}
	*l = evts
	return nil
}

func UnmarshalJsonEvents(data []byte) ([]Event, error) {
	var gEvts []marshalStruct
	err := json.Unmarshal(data, &gEvts)
	if err != nil {
		return nil, xerrors.Errorf("unmarshal evts: %w", err)
	}

	evts := make([]Event, 0, len(gEvts))
	for _, gEvt := range gEvts {
		var e Event
		switch gEvt.EventType {
		case LeftClickEvent:
			e = &ClickEvent{}
		case PlayerJoined:
			e = &PlayerJoin{}
		}
		err := json.Unmarshal(gEvt.Payload, e)
		if err != nil {
			return nil, xerrors.Errorf("unmarshal evt: %w", err)
		}
		evts = append(evts, e)
	}

	return evts, err
}

func MarshalJsonEvents(evts []Event) ([]byte, error) {
	gEvts := make([]marshalStruct, len(evts))
	for i, evt := range evts {
		payload, err := json.Marshal(evt)
		if err != nil {
			return nil, xerrors.Errorf("marshal: %w", err)
		}
		gEvts[i] = marshalStruct{
			EventType: evt.Type(),
			Payload:   payload,
		}
	}

	data, err := json.Marshal(gEvts)
	if err != nil {
		return nil, xerrors.Errorf("marshal gEvts: %w", err)
	}
	return data, nil
}

type Event interface {
	GetID() uint64
	Type() EventType
	// Tick will allow the event to advance 1 tick. If the event is done, it should return a nil.
	Tick(w *world2.World) (Event, error)

	// Should optimize data later
	//MarshalBinary() ([]byte, error)
	//UnmarshalBinary(data []byte) error
	//UnmarshalBinaryData(data []byte) ([]byte, error)
}

func AddLogFields(l *zerolog.Event, e Event) *zerolog.Event {
	return l.Uint64("ID", e.GetID()).Str("type", string(e.Type())).Str("log_type", "game_event")
}
