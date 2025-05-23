package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/csotherden/prescription-parser/pkg/config"
	"github.com/csotherden/prescription-parser/pkg/datastore"
	"github.com/csotherden/prescription-parser/pkg/handlers/parser"
	parserPkg "github.com/csotherden/prescription-parser/pkg/parser"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	config config.Config
	logger *zap.Logger
	router *mux.Router
	server *http.Server
	ds     datastore.Datastore
	parser parserPkg.Parser
}

func NewServer(cfg config.Config, ds datastore.Datastore, parser parserPkg.Parser, logger *zap.Logger) (*Server, error) {
	s := &Server{
		config: cfg,
		logger: logger,
		router: mux.NewRouter(),
		ds:     ds,
		parser: parser,
	}

	s.setupMiddleware()
	if err := s.setupRoutes(); err != nil {
		return nil, fmt.Errorf("failed to setup routes: %w", err)
	}

	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:      s.router,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return s, nil
}

func (s *Server) setupMiddleware() {
	s.router.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(handlePreflight)
	s.router.Use(corsMiddleware)
}

func (s *Server) setupRoutes() error {
	// Create handlers
	parserHandler := parser.NewHandler(s.parser, s.ds, s.logger)

	// Setup API routes
	apiRouter := s.router.PathPrefix("/api").Subrouter()

	// Health check
	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is running"))
	}).Methods("GET")

	// Register the API routes
	parserHandler.RegisterRoutes(apiRouter)

	return nil
}

func (s *Server) Start() error {
	s.logger.Info("Starting server", zap.String("address", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Middleware functions
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Project-Id")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handlePreflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, HEAD, OPTIONS, PATCH, POST, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(http.StatusOK)
}
