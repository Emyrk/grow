package network

import (
	"context"
	"encoding/json"

	"github.com/emyrk/grow/game"
	"github.com/emyrk/grow/game/events"
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
			case message.MTEventSync:
				var ms message.EventSync
				err := json.Unmarshal(msg.Payload, &ms)
				if err != nil {
					gc.Log.Err(err).Msg("unmarshal sync")
					continue
				}

				evts, err := events.UnmarshalJsonEvents(ms.Events)
				if err != nil {
					gc.Log.Err(err).Msg("unmarshal evts")
					continue
				}

				gc.Log.Info().Uint64("tick", ms.GameTick).Int("event_count", len(evts)).Msg("event sync")
				gc.ReceiveGameEvents(ctx, ms.GameTick, evts)
			}
		}
	}
}
