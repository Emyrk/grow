package game

import (
	"sync"

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
}

func NewGameServer(log zerolog.Logger, cfg GameConfig) *GameServer {
	c := &GameServer{
		G:         NewGame(log, cfg),
		Listeners: make(map[world.PlayerID]*ListeningPlayer),
		log:       log,
	}
	return c
}

func (g *GameServer) Update() error {
	syncEvts, evts := g.G.Update(g.Gametick)
	if syncEvts {
		// Broadcast events
		g.log.Info().Uint64("gametick", g.Gametick).Int("new_events_count", len(evts)).Msg("tick")
		g.mu.RLock()
		// TODO: Might want to do this in another go routine. Currently this can slow down the server.
		for _, l := range g.Listeners {
			l.Broadcast(g.Gametick, evts)
		}
		g.mu.RUnlock()
	}
	g.Gametick++

	return nil
}

func (c *GameServer) SendEvents(es []events.Event) {
	for i := range es {
		err := c.G.EC.SendEvent(es[i])
		if err != nil {
			c.log.Err(err).Msg("send event")
		}
	}
}

// AddListener adds a new listener to the game. The listener can submit event requests through the events
// channel.
func (c *GameServer) AddListener(id world.PlayerID, broadcast ProcessEvents) chan<- events.Event {
	c.mu.Lock()
	defer c.mu.Unlock()

	p := &ListeningPlayer{
		ID:        id,
		NewEvents: make(chan events.Event, 100),
		Broadcast: broadcast,
	}
	c.Listeners[id] = p
	return p.NewEvents
}

func (c *GameServer) RemoveListener(id world.PlayerID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	p, ok := c.Listeners[id]
	if ok {
		close(p.NewEvents)
	}
}

type ListeningPlayer struct {
	ID world.PlayerID
	// NewEvents are events the player has submitted to be recorded in the controller.
	NewEvents chan events.Event
	Broadcast func(gametick uint64, events []events.Event)
}

// WatchPlayer will watch the player for any events they wish to submit to the event controller.
func (c *GameServer) WatchPlayer(p *ListeningPlayer) {
	for {
		select {
		case e, ok := <-p.NewEvents:
			if !ok {
				return // Player has left
			}

			err := c.G.EC.SendEvent(e)
			if err != nil {
				// TODO: Notify the user
				c.log.Err(err).Uint16("pid", uint16(p.ID)).Msg("send player event")
				break
			}
		}
	}
}
