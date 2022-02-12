package server

import (
	"fmt"
	"net/http"

	"github.com/emyrk/grow/game"

	"github.com/rs/zerolog"

	"github.com/go-chi/chi/v5"
)

// Webserver can run a single game.
type Webserver struct {
	Srv  *http.Server
	Log  zerolog.Logger
	Game *game.GameServer
}

type ServerConfig struct {
	Port int
	Log  zerolog.Logger
}

func NewWebserver(cfg ServerConfig, gm *game.GameServer) *Webserver {
	mux := chi.NewRouter()
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: mux,
	}

	gs := &Webserver{
		Srv:  srv,
		Log:  cfg.Log.With().Str("server", "server").Logger(),
		Game: gm,
	}

	AttachRoutes(gs, mux)

	return gs
}

func (gs *Webserver) Start() error {
	gs.Log.Info().Str("addr", gs.Srv.Addr).Msg("server listening")
	return gs.Srv.ListenAndServe()
}

func (gs *Webserver) Close() error {
	return gs.Srv.Close()
}
