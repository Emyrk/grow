package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/emyrk/grow/game/events"
	"github.com/emyrk/grow/internal/crand"
	"github.com/emyrk/grow/server/message"
	"github.com/emyrk/grow/world"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// HandleGame handles the communication between the server and clients for the game state.
func (gs *Webserver) HandleGame(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		_ = json.NewEncoder(w).Encode(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := world.PlayerID(crand.Uint64())
	stopClient := func() {
		cancel()
		gs.Game.RemoveListener(id)
		c.Close(websocket.StatusInternalError, "connection closed")
	}
	defer stopClient()

	gs.Game.AddListener(id, func(gametick uint64, evts []events.Event) {
		// Broadcast to the player
		data, err := events.MarshalJsonEvents(evts)
		if err != nil {
			gs.Log.Err(err).Msg("marshal events")
			return
		}

		msg := message.EventSync{
			GameTick: gametick,
			Events:   data,
		}
		data, err = json.Marshal(msg)
		if err != nil {
			gs.Log.Err(err).Msg("marshal event sync")
			return
		}

		data, err = json.Marshal(message.SocketMessage{
			MessageType: message.MTEventSync,
			Payload:     data,
		})
		if err != nil {
			gs.Log.Err(err).Msg("marshal socket msg")
			return
		}

		err = c.Write(ctx, websocket.MessageText, data)
		if err != nil {
			gs.Log.Err(err).Msg("player broadcast")
			stopClient()
			return
		}
	})

	log := gs.Log.With().Uint16("pid", uint16(id)).Logger()
	for {
		select {
		case <-ctx.Done():
			stopClient()
			return
		default:

		}
		var msg message.SocketMessage
		err := wsjson.Read(ctx, c, &msg)
		//_, data, err := c.Read(ctx)
		if err != nil {
			log.Err(err).Msg("websocket read")
			break
		}

		switch msg.MessageType {
		case message.MTNewEvents:
			evts, err := events.UnmarshalJsonEvents(msg.Payload)
			if err != nil {
				log.Warn().Uint64("msg_type", msg.MessageType).Str("payload", string(msg.Payload)).Msg("unmarshal event")
				break
			}
			gs.Game.SendEvents(evts)
		default:
			log.Warn().Uint64("msg_type", msg.MessageType).Msg("unknown type")
		}
		log.Debug().Uint64("msg_type", msg.MessageType).Msg("message from player")
	}
}
