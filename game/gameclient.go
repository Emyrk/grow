package game

import (
	"context"
	"encoding/json"
	"sync"

	"golang.org/x/xerrors"

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
	// SendGameMessage is for pushing events to a game server.
	SendGameMessage func(msgType GameMessageType, payload []byte) error
	ClientMode      bool
}

func NewGameClient(log zerolog.Logger, cfg GameConfig) *GameClient {
	c := &GameClient{
		G:           NewGame(log, cfg),
		Log:         log,
		syncedTicks: make(map[uint64][]events.Event),
	}
	c.SendGameMessage = c.GameMessage
	return c
}

func (g *GameClient) GameMessage(msgType GameMessageType, data []byte) error {
	switch msgType {
	case MsgGameNewEvents:
		// Only local mode should handle this
		if !g.ClientMode {
			var msg NewEvents
			err := json.Unmarshal(data, &msg)
			if err != nil {
				return xerrors.Errorf("unmarshal game sync: %w", err)
			}

			for i := range msg.Eventlist {
				err := g.G.EC.SendEvent(msg.Eventlist[i])
				if err != nil {
					g.Log.Err(err).Msg("send evt")
				}
			}

		}
	case MsgGameSync:
		var msg GameSync
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return xerrors.Errorf("unmarshal game sync: %w", err)
		}

		// Lets do a full sync
		if msg.GameTick > g.Gametick {
			g.fullSync(msg)
		}
	case MsgTickEventList:
		var msg TickEventList
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return xerrors.Errorf("unmarshal tick event list: %w", err)
		}
		g.ReceiveGameEvents(msg.GameTick, msg.Eventlist)
		g.Log.Info().Uint64("tick", msg.GameTick).Uint64("behind", msg.GameTick-g.Gametick).Int("event_count", len(msg.Eventlist)).Msg("event sync")
	default:
		return xerrors.Errorf("msg type %s not recognized", msgType)
	}
	return nil
}

func (c *GameClient) fullSync(gameSync GameSync) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.G.World = gameSync.World
	c.G.EC.ReplaceEventList(gameSync.EventList)
	c.Gametick = gameSync.GameTick
}

// ReceiveGameEvents allows us to make advancements in our ticks
func (c *GameClient) ReceiveGameEvents(gametick uint64, evts []events.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if gametick < c.Gametick {
		log.Warn().Uint64("msg_tick", gametick).Uint64("tick", c.Gametick).Msg("server sent us a tick we are past")
		return // we are past that tick
	}

	c.syncedTicks[gametick] = evts
}

func (c *GameClient) UseServer(send func(msgType GameMessageType, payload []byte) error) *GameClient {
	c.SendGameMessage = send
	c.ClientMode = true
	return c
}

func (c *GameClient) Update() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	eventPoint := events.SyncTick(c.Gametick)
	if eventPoint && c.ClientMode {
		evts, ok := c.syncedTicks[c.Gametick]
		if !ok {
			c.waiting++
			if c.waiting%60 == 0 {
				// TODO: Request a state sync to get past this
				c.Log.Warn().Uint64("tick", c.Gametick).Msg("waiting for sync")
				c.waiting = 0
				err := c.SendGameMessage(MsgGameSync, []byte("{}"))
				if err != nil {
					return xerrors.Errorf("request game sync: %w", err)
				}
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
