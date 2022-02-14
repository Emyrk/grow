package game

import (
	world2 "github.com/emyrk/grow/game/world"
)


// CreateGameSync creates a full game sync payload to bring a client up to speed without them having to replay from
// game tick 0.
type CreateGameSync struct {
	// Players are the players who need to be synced
	Players []world2.PlayerID
}


