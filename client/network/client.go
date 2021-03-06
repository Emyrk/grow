package network

import (
	"context"
	"encoding/json"
	"io"

	"github.com/emyrk/grow/internal/network"

	"github.com/emyrk/grow/game"

	"nhooyr.io/websocket/wsjson"

	"github.com/emyrk/grow/server/message"
	"github.com/rs/zerolog"
	"golang.org/x/xerrors"
	"nhooyr.io/websocket"
)

type NetworkClient struct {
	Conn *websocket.Conn
	Log  zerolog.Logger
}

func Connect(ctx context.Context, log zerolog.Logger, address string) (*NetworkClient, error) {
	c, _, err := websocket.Dial(ctx, address, &websocket.DialOptions{})
	if err != nil {
		return nil, xerrors.Errorf("dial: %w", err)
	}
	c.SetReadLimit(network.ReadLimit)

	return &NetworkClient{
		Conn: c,
		Log:  log,
	}, nil
}

func (c *NetworkClient) Close() error {
	return c.Conn.Close(websocket.StatusNormalClosure, "leaving")
}

func (c *NetworkClient) ReadMessages(ctx context.Context) <-chan *message.SocketMessage {
	msgs := make(chan *message.SocketMessage, 100)
	go func() {
		for {
			var msg message.SocketMessage
			err := wsjson.Read(ctx, c.Conn, &msg)
			//_, data, err := c.Conn.Read(ctx)
			if err != nil {
				if xerrors.Is(err, io.EOF) {
					continue
				}
				c.Log.Err(err).Msg("read msg")
				close(msgs)
				return
			}

			//var msg message.SocketMessage
			//err = json.Unmarshal(data, &msg)
			//if err != nil {
			//	c.Log.Err(err).Msg("unmarshal msg")
			//	continue
			//}

			msgs <- &msg
		}
	}()
	return msgs
}

func (c *NetworkClient) SendGameMessage(ctx context.Context) func(msgType game.GameMessageType, payload []byte) error {
	return func(msgType game.GameMessageType, payload []byte) error {
		smsg, err := json.Marshal(message.SocketMessage{
			MessageType: message.MTGameMessage,
			PayloadType: msgType,
			Payload:     payload,
		})
		if err != nil {
			return xerrors.Errorf("marshal socket msg: %w", err)
		}
		err = c.Conn.Write(ctx, websocket.MessageText, smsg)
		if err != nil {
			return xerrors.Errorf("send msg: %w", err)
		}
		return nil
	}
}
