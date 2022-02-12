package server

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/emyrk/grow/game/events"
	"github.com/emyrk/grow/server/message"
	"github.com/emyrk/grow/world"
	"nhooyr.io/websocket"
)

// HandleGame handles the communication between the server and clients for the game state.
func (gs *GameServer) HandleGame(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		_ = json.NewEncoder(w).Encode(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer c.Close(websocket.StatusInternalError, "connection closed")

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	id := world.PlayerID(rand.Uint64())
	defer gs.Game.RemoveListener(id)
	gs.Game.AddListener(id, func(gametick uint64, evts []events.Event) {
		// Broadcast to the player
		data, err := events.MarshalJsonEvents(evts)
		if err != nil {
			cancel()
			gs.Log.Err(err).Msg("marshal events")
			return
		}

		msg := message.EventSync{
			GameTick: gametick,
			Events:   data,
		}
		data, err = json.Marshal(msg)
		if err != nil {
			cancel()
			gs.Log.Err(err).Msg("marshal event sync")
			return
		}

		err = c.Write(ctx, websocket.MessageText, data)
		if err != nil {
			cancel()
			gs.Log.Err(err).Msg("player broadcast")
			return
		}
	})

	log := gs.Log.With().Uint16("pid", uint16(id)).Logger()
	for {
		_, data, err := c.Reader(ctx)
		if err != nil {
			log.Err(err).Msg("websocket read")
			break
		}

		var msg message.SocketMessage
		err = json.NewDecoder(data).Decode(&msg)
		if err != nil {
			log.Err(err).Msg("decode read")
			break
		}

		switch msg.MessageType {
		case message.MTNewEvents:
			evts, err := events.UnmarshalJsonEvents(msg.Payload)
			if err != nil {
				log.Warn().Uint64("msg_type", msg.MessageType).Msg("unmarshal event")
				break
			}
			gs.Game.SendEvents(evts)
		default:
			log.Warn().Uint64("msg_type", msg.MessageType).Msg("unknown type")
		}
	}

	_ = c.Close(websocket.StatusNormalClosure, "")
}
