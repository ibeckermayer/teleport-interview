package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ibeckermayer/teleport-interview/backend/internal/handlers"
)

// Config is the top level config object.
type Config struct {
	port         int    // -port; default 8000
	certFilePath string // -cert; default "../certs/localhost.crt"
	keyFilePath  string // -key ; default "../certs/localhost.key"
}

// NewConfig creates a new server.Config
func NewConfig(port int, certFilePath string, keyFilePath string) Config {
	return Config{port, certFilePath, keyFilePath}
}

// Server object initializes route handlers and external connections, and serves application
type Server struct {
	cfg    Config
	router *mux.Router
	// TODO: database and session manager connections will go here
}

// New initializes routes and handlers and returns a ready-to-run server
func New(cfg Config) *Server {
	srv := &Server{cfg, mux.NewRouter()}

	loginHandler := handlers.NewLoginHandler()
	srv.router.Handle("/api/login", loginHandler).Methods("POST")

	// NOTE: It's important that this handler be registered after the other handlers, or else
	// all routes return a 404 (at least in development). TODO: figure out why this is the case.
	spaHandler := handlers.NewSpaHandler("../frontend", "index.html")
	srv.router.PathPrefix("/").Handler(spaHandler)

	return srv
}

// Run starts the server
func (srv *Server) Run() error {
	log.Printf("Server listening on port %v", srv.cfg.port)
	return http.ListenAndServeTLS(fmt.Sprintf(":%v", srv.cfg.port), srv.cfg.certFilePath, srv.cfg.keyFilePath, srv.router)
}
