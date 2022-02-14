package server

import (
	"context"
	"encoding/json"
	world2 "github.com/emyrk/grow/game/world"
	"net/http"

	"github.com/emyrk/grow/internal/network"

	"github.com/emyrk/grow/game"

	"github.com/emyrk/grow/internal/crand"
	"github.com/emyrk/grow/server/message"
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
	c.SetReadLimit(network.ReadLimit)

	id := world2.PlayerID(crand.Uint64())
	stopClient := func() {
		cancel()
		gs.Game.RemoveListener(id)
		c.Close(websocket.StatusInternalError, "connection closed")
	}
	defer stopClient()

	gs.Game.AddListener(id, func(msgType game.GameMessageType, data []byte) {
		data, err = json.Marshal(message.SocketMessage{
			MessageType: message.MTGameMessage,
			PayloadType: msgType,
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

	log := gs.Log.With().Uint64("pid", id).Logger()
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
		case message.MTGameMessage:
			// Forward all game messages to the game
			err := gs.Game.GameMessage(id, msg.PayloadType, msg.Payload)
			if err != nil {
				log.Warn().Uint64("msg_type", msg.MessageType).Str("pay_type", msg.PayloadType).Msg("game msg")
				break
			}
		default:
			log.Warn().Uint64("msg_type", msg.MessageType).Msg("unknown type")
		}
		log.Debug().Uint64("msg_type", msg.MessageType).Msg("message from player")
	}
}
