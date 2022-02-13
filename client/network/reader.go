package network

import (
	"context"

	"github.com/emyrk/grow/game"
	"github.com/emyrk/grow/server/message"
)

func HandleSocketMessages(ctx context.Context, gc *game.GameClient, msgs <-chan *message.SocketMessage) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}
			switch msg.MessageType {
			case message.MTGameMessage:
				err := gc.GameMessage(msg.PayloadType, msg.Payload)
				if err != nil {
					gc.Log.Err(err).Str("type", msg.PayloadType).Msg("game msg")
					continue
				}
			}
		}
	}
}
