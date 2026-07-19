package api

import (
	"net/http"
	"sync/atomic"

	"github.com/Robot-tim1/Chirpy/internal/database"
)

type Server struct {
	mux *http.ServeMux
	cfg *apiConfig
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
}

type ApiConfigParams struct {
	Db       database.DBTX
	Platform string
	Secret   string
}

func NewServer(params ApiConfigParams) *Server {
	s := &Server{
		mux: http.NewServeMux(),
		cfg: &apiConfig{
			fileserverHits: atomic.Int32{},
			db:             database.New(params.Db),
			platform:       params.Platform,
			secret:         params.Secret,
		},
	}
	s.registerRoutes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) registerRoutes() {
	fsHandler := s.cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	s.mux.Handle("/app/", fsHandler)

	s.mux.HandleFunc("GET /admin/metrics", s.cfg.handlerRequestNum)
	s.mux.HandleFunc("POST /admin/reset", s.cfg.handlerResetEnd)
	s.mux.HandleFunc("GET /api/healthz", handlerHealthzEnd)
	s.mux.HandleFunc("POST /api/chirps", s.cfg.handlerChirpPost)
	s.mux.HandleFunc("GET /api/chirps", s.cfg.handlerChirpGet)
	s.mux.HandleFunc("GET /api/chirps/{chirpID}", s.cfg.handlerChirpGetID)
	s.mux.HandleFunc("DELETE /api/chirps/{chirpID}", s.cfg.handlerChirpDeleteID)
	s.mux.HandleFunc("POST /api/users", s.cfg.handlerUserPost)
	s.mux.HandleFunc("PUT /api/users", s.cfg.handlerUserPut)
	s.mux.HandleFunc("POST /api/login", s.cfg.handlerLoginEnd)
	s.mux.HandleFunc("POST /api/refresh", s.cfg.handlerRefreshEnd)
	s.mux.HandleFunc("POST /api/revoke", s.cfg.handlerRevokeEnd)
}
