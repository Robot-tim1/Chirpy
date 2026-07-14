package api

import (
	"net/http"
	"sync/atomic"
)

type Server struct {
	mux       *http.ServeMux
	apiConfig *apiConfig
	// db   *database.DB
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func NewServer() *Server {
	s := &Server{
		mux:       http.NewServeMux(),
		apiConfig: &apiConfig{atomic.Int32{}},
	}
	s.registerRoutes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) registerRoutes() {
	fsHandler := s.apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	s.mux.Handle("/app/", fsHandler)

	s.mux.HandleFunc("GET /admin/metrics", s.apiConfig.handlerRequestNum)
	s.mux.HandleFunc("POST /admin/reset", s.apiConfig.handlerResetNum)
	s.mux.HandleFunc("GET /api/healthz", handlerHealthzEnd)
	s.mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirpEnd)

}
