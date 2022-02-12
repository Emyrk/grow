package server

import "github.com/go-chi/chi/v5"

func AttachRoutes(gs *Webserver, rs chi.Router) {
	rs.Get("/version", gs.VersionHandler)
	rs.HandleFunc("/socket", gs.HandleGame)
}
