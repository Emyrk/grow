package server

import "github.com/go-chi/chi/v5"

func AttachRoutes(gs *GameServer, rs chi.Router) {
	rs.Get("/version", gs.VersionHandler)
}
