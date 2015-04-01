package main

import (
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/thoas/stats"
	"encoding/json"
	"github.com/zenazn/goji"
	"net/http"
)

func main() {
	middleware := stats.New()

	goji.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})
	goji.Get("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		stats := middleware.Data()

		b, _ := json.Marshal(stats)

		w.Write(b)
	})

	goji.Use(middleware.Handler)
	goji.Serve()
}
