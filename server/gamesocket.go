package server

import (
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/xerrors"

	"nhooyr.io/websocket/wsjson"

	"github.com/emyrk/grow/game/events"
	"github.com/emyrk/grow/server/message"
	"github.com/emyrk/grow/world"
	"nhooyr.io/websocket"
)

// HandleGame handles the communication between the server and clients for the game state.
func (gs *Webserver) HandleGame(w http.ResponseWriter, r *http.Request) {
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

		data, err = json.Marshal(message.SocketMessage{
			MessageType: message.MTEventSync,
			Payload:     data,
		})
		if err != nil {
			cancel()
			gs.Log.Err(err).Msg("marshal socket msg")
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
		var msg message.SocketMessage
		err := wsjson.Read(ctx, c, &msg)
		//_, data, err := c.Read(ctx)
		if err != nil {
			if xerrors.Is(err, io.EOF) {
				continue
			}
			log.Err(err).Msg("websocket read")
			break
		}

		//var msg message.SocketMessage
		//err = json.Unmarshal(data, &msg)
		//if err != nil {
		//	log.Err(err).Msg("decode read")
		//	break
		//}

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

	_ = c.Close(websocket.StatusNormalClosure, "")
}
