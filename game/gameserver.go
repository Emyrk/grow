package game

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"golang.org/x/xerrors"

	"github.com/emyrk/grow/game/events"
	"github.com/emyrk/grow/world"
	"github.com/rs/zerolog"
)

// GameServer handles managing all the players listening to a game.
// A listening player can be a spectator.
type GameServer struct {
	G        *Game
	Gametick uint64

	log       zerolog.Logger
	mu        sync.RWMutex
	Listeners map[world.PlayerID]*ListeningPlayer
	Commands  chan interface{}
}

func NewGameServer(log zerolog.Logger, cfg GameConfig) *GameServer {
	c := &GameServer{
		G:         NewGame(log, cfg),
		Listeners: make(map[world.PlayerID]*ListeningPlayer),
		log:       log,
		Commands:  make(chan interface{}, 100),
	}
	return c
}

func (g *GameServer) GameLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := g.Update()
			if err != nil {
				g.log.Err(err).Msg("on tick")
			}
		}
	}
}

func (g *GameServer) GameMessage(playerID world.PlayerID, msgType GameMessageType, data []byte) error {
	log := g.log.With().Uint16("pid", playerID).Str("msg_type", msgType).Int("pay_size", len(data)).Logger()
	switch msgType {
	case MsgTickEventList:
		// TODO: Player needs these events
		fallthrough // For now, just send them a full sync
	case MsgGameSync:
		// Send the player a sync
		cmd := &CreateGameSync{Players: []world.PlayerID{playerID}}
		g.Commands <- cmd
		log.Info().Msg("request game sync")
	case MsgGameNewEvents:
		var msg NewEvents
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return xerrors.Errorf("unmarshal new evts: %w", err)
		}
		g.SendEvents(msg.Eventlist)
		log.Info().Int("evt_cnt", len(msg.Eventlist)).Msg("new events")
	default:
		return xerrors.Errorf("msg type %s not recognized", msgType)
	}
	return nil
}

func (g *GameServer) Update() error {
	syncEvts, evts := g.G.Update(g.Gametick)
	if syncEvts {
		// Broadcast events
		g.log.Info().Uint64("gametick", g.Gametick).Int("new_events_count", len(evts)).Msg("tick")
		g.mu.RLock()

		msg := TickEventList{
			GameTick:  g.Gametick,
			Eventlist: evts,
		}
		data, err := json.Marshal(msg)
		if err != nil {
			return xerrors.Errorf("marshal tick events: %w", err)
		}

		// TODO: Might want to do this in another go routine. Currently this can slow down the server.
		for _, l := range g.Listeners {
			l.BroadcastData(msg.Type(), data)
		}
		g.mu.RUnlock()
	}

	g.handleCommands()

	g.Gametick++

	return nil
}

func (g *GameServer) handleCommands() {
	syncPlayers := make(map[world.PlayerID]struct{})

CmdLoop:
	for {
		select {
		case cmd := <-g.Commands:
			switch cmd.(type) {
			case *CreateGameSync:
				gs, ok := cmd.(*CreateGameSync)
				if ok {
					for _, p := range gs.Players {
						syncPlayers[p] = struct{}{}
					}
				}
			default:
				g.log.Error().Str("type", reflect.TypeOf(cmd).String()).Msg("unknown command")
			}
		default:
			break CmdLoop
		}
	}

	// Handle sync events
	if len(syncPlayers) > 0 {
		state := GameSync{
			World:     g.G.World,
			EventList: g.G.EC.EventList(),
			GameTick:  g.Gametick,
		}
		worldData, err := json.Marshal(state)
		if err != nil {
			g.log.Err(err).Msg("marshal game state")
		} else {
			for k := range syncPlayers {
				g.log.Info().Uint16("pid", k).Int("pay_size", len(worldData)).Msg("send game state")
				g.Listeners[k].BroadcastData(state.Type(), worldData)
			}
		}
	}
}

func (g *GameServer) SendEvents(es []events.Event) {
	for i := range es {
		err := g.G.EC.SendEvent(es[i])
		if err != nil {
			g.log.Err(err).Msg("send event")
		}
	}
}

// AddListener adds a new listener to the game. The listener can submit event requests through the events
// channel.
func (g *GameServer) AddListener(id world.PlayerID, broadcast BroadcastGameMessage) chan<- events.Event {
	g.mu.Lock()
	defer g.mu.Unlock()

	p := &ListeningPlayer{
		ID:            id,
		NewEvents:     make(chan events.Event, 100),
		BroadcastData: broadcast,
	}
	g.Listeners[id] = p
	return p.NewEvents
}

func (g *GameServer) RemoveListener(id world.PlayerID) {
	g.mu.Lock()
	defer g.mu.Unlock()

	p, ok := g.Listeners[id]
	delete(g.Listeners, id)
	if ok {
		close(p.NewEvents)
	}
}

type ListeningPlayer struct {
	ID world.PlayerID
	// NewEvents are events the player has submitted to be recorded in the controller.
	NewEvents     chan events.Event
	BroadcastData func(msgType GameMessageType, data []byte)
}

// WatchPlayer will watch the player for any events they wish to submit to the event controller.
func (g *GameServer) WatchPlayer(p *ListeningPlayer) {
	for {
		select {
		case e, ok := <-p.NewEvents:
			if !ok {
				return // Player has left
			}

			err := g.G.EC.SendEvent(e)
			if err != nil {
				// TODO: Notify the user
				g.log.Err(err).Uint16("pid", uint16(p.ID)).Msg("send player event")
				break
			}
		}
	}
}
