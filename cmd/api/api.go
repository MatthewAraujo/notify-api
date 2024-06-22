package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/MatthewAraujo/notify/service/events"
	"github.com/MatthewAraujo/notify/service/health"
	"github.com/MatthewAraujo/notify/service/notifications"
	"github.com/MatthewAraujo/notify/service/repository"
	"github.com/MatthewAraujo/notify/service/user"
	"github.com/MatthewAraujo/notify/service/webhooks"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Start() error {
	router := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow only the specific origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Allow specific headers
		AllowCredentials: true,
		Debug:            true, // Enable debug mode for CORS
	})

	router.Use(c.Handler)

	// if the api changes in the future we can just change the version here, and the old version will still be available
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// create a new handler
	healthHandler := health.NewHandler()
	// register the handler
	healthHandler.Register(subrouter)

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.Register(subrouter)

	notificationStore := notifications.NewStore(s.db)
	notificationHandler := notifications.NewHandler(notificationStore)
	notificationHandler.Register(subrouter)

	webhooksStore := webhooks.NewStore(s.db)
	webhooksHandler := webhooks.NewHandler(webhooksStore)
	webhooksHandler.Register(subrouter)

	repositoryStore := repository.NewStore(s.db)
	repositoryHandler := repository.NewHandler(repositoryStore)
	repositoryHandler.Register(subrouter)

	eventTypeStore := events.NewStore(s.db)
	eventTypeHandler := events.NewHandler(eventTypeStore)
	eventTypeHandler.Register(subrouter)

	log.Println("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
