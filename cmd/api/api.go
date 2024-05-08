package api

import (
	"log"
	"net/http"

	"github.com/MatthewAraujo/notify/cmd/health"
	"github.com/MatthewAraujo/notify/cmd/service/webhooks"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{addr: addr}
}

func (s *APIServer) Start() error {
	router := mux.NewRouter()
	// if the api changes in the future we can just change the version here, and the old version will still be available
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// create a new handler
	healthHandler := health.NewHandler()
	// register the handler
	healthHandler.Register(subrouter)

	webhooksHandler := webhooks.NewHandler()
	webhooksHandler.Register(subrouter)
	log.Println("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
