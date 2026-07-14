package main

import (
	"net/http"

	"github.com/Robot-tim1/Chirpy/internal/api"
)

func main() {
	server := api.NewServer()

	httpServer := &http.Server{
		Handler: server.Handler(),
		Addr:    ":8080",
	}
	httpServer.ListenAndServe()
}
