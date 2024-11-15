package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", Chain(s.ForwardHandler, Cache(s), Logging()))
	// mux.HandleFunc("/ping", s.Ping)

	// mux.HandleFunc("/health", s.healthHandler)

	return mux
}

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("ok")
}

func (s *Server) ForwardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Cache", "MISS")
	proxy := createProxy(s.target, s, r)
	proxy.ServeHTTP(w, r)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
