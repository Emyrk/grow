package server

import (
	"encoding/json"
	"net/http"

	"github.com/emyrk/grow/internal/version"
)

type VersionStruct struct {
	Version      string `json:"version"`
	CommitSHA1   string `json:"commit_sha_1"`
	CompiledDate string `json:"compiled_date"`
}

func (GameServer) VersionHandler(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(VersionStruct{
		Version:      version.Version,
		CommitSHA1:   version.CommitSHA1,
		CompiledDate: version.CompiledDate,
	})
}
