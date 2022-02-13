package game

import "github.com/emyrk/grow/world"


// CreateGameSync creates a full game sync payload to bring a client up to speed without them having to replay from
// game tick 0.
type CreateGameSync struct {
	// Players are the players who need to be synced
	Players []world.PlayerID
}


