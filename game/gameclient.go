package game

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/emyrk/grow/game/events"
	"github.com/rs/zerolog"
)

type ClientSendEvents func(ctx context.Context, evts []events.Event) error

// GameClient is able to run a game, but relies on a server to handle event ordering.
type GameClient struct {
	G           *Game
	mu          sync.Mutex
	Gametick    uint64
	waiting     uint64
	syncedTicks map[uint64][]events.Event

	Log zerolog.Logger
	// SendGameEvents is for pushing events to a game server.
	SendGameEvents ClientSendEvents
}

func NewGameClient(log zerolog.Logger, cfg GameConfig) *GameClient {
	c := &GameClient{
		G:           NewGame(log, cfg),
		Log:         log,
		syncedTicks: make(map[uint64][]events.Event),
	}
	if c.SendGameEvents == nil {
		c.SendGameEvents = func(_ context.Context, evts []events.Event) error {
			for i := range evts {
				err := c.G.EC.SendEvent(evts[i])
				if err != nil {
					c.Log.Err(err).Msg("send evt")
				}
			}
			return nil
		}
	}
	return c
}

// ReceiveGameEvents allows us to make advancements in our ticks
func (c *GameClient) ReceiveGameEvents(ctx context.Context, gametick uint64, evts []events.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if gametick < c.Gametick {
		log.Warn().Uint64("msg_tick", gametick).Uint64("tick", c.Gametick).Msg("server sent us a tick we are past")
		return // we are past that tick
	}

	c.syncedTicks[gametick] = evts
}

func (c *GameClient) UseServer(send ClientSendEvents) *GameClient {
	c.SendGameEvents = send
	return c
}

func (c *GameClient) Update() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	eventPoint := events.SyncTick(c.Gametick)
	if eventPoint {
		evts, ok := c.syncedTicks[c.Gametick]
		if !ok {
			c.waiting++
			if c.waiting%60 == 0 {
				// TODO: Request a state sync to get past this
				c.Log.Warn().Uint64("tick", c.Gametick).Msg("waiting for sync")
				c.waiting = 0
			}
			return nil
		}

		// Send all events in the order to be processed
		for i := range evts {
			err := c.G.EC.SendEvent(evts[i])
			if err != nil {
				c.Log.Err(err).Msg("GAME WILL DESYNC, SERVER EVENT REJECTED")
			}
		}
		delete(c.syncedTicks, c.Gametick)
	}

	// TODO: Handle blocking for events
	c.G.Update(c.Gametick)
	c.Gametick++
	return nil
}
