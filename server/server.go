package server

import (
	"fmt"
	"net/http"

	"github.com/emyrk/grow/game"

	"github.com/rs/zerolog"

	"github.com/go-chi/chi/v5"
)

// GameServer can run a single game.
type GameServer struct {
	Srv  *http.Server
	Log  zerolog.Logger
	Game *game.GameServer
}

type ServerConfig struct {
	Port int
	Log  zerolog.Logger
}

func NewGameServer(cfg ServerConfig, gm *game.GameServer) *GameServer {
	mux := chi.NewRouter()
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: mux,
	}

	gs := &GameServer{
		Srv:  srv,
		Log:  cfg.Log.With().Str("server", "server").Logger(),
		Game: gm,
	}

	AttachRoutes(gs, mux)

	return gs
}

func (gs *GameServer) Start() error {
	gs.Log.Info().Str("addr", gs.Srv.Addr).Msg("server listening")
	return gs.Srv.ListenAndServe()
}

func (gs *GameServer) Close() error {
	return gs.Srv.Close()
}
