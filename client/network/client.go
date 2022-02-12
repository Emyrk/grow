package network

import (
	"context"

	"golang.org/x/xerrors"
	"nhooyr.io/websocket"
)

type NetworkClient struct {
	Conn *websocket.Conn
}

func Connect(ctx context.Context, address string) (*NetworkClient, error) {
	c, _, err := websocket.Dial(ctx, address, &websocket.DialOptions{})
	if err != nil {
		return nil, xerrors.Errorf("dial: %w", err)
	}

	return &NetworkClient{
		Conn: c,
	}, nil
}

func (c *NetworkClient) Close() error {
	return c.Conn.Close(websocket.StatusNormalClosure, "leaving")
}
