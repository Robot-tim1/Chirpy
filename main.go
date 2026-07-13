package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiConfig := apiConfig{atomic.Int32{}}

	serveMux := http.NewServeMux()
	fsHandler := apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	serveMux.Handle("/app/", fsHandler)

	serveMux.HandleFunc("GET /admin/metrics", apiConfig.handlerRequestNum)
	serveMux.HandleFunc("POST /admin/reset", apiConfig.handlerResetNum)
	serveMux.HandleFunc("GET /api/healthz", handlerEndpoint)

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
