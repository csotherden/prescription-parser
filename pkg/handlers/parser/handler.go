package parser

import (
	"fmt"

	"github.com/csotherden/prescription-parser/pkg/config"
	"github.com/csotherden/prescription-parser/pkg/datastore"
	"github.com/csotherden/prescription-parser/pkg/parser"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Handler handles assistant-related API requests
type Handler struct {
	ds     *datastore.Datastore
	logger *zap.Logger
	parser *parser.PrescriptionParser
}

// NewHandler creates a new assistant handler instance
func NewHandler(cfg config.Config, ds *datastore.Datastore, logger *zap.Logger) (*Handler, error) {
	// Create a new parser
	parserInstance, err := parser.NewParser(cfg, ds, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create parser handler. could not create parser: %w", err)
	}

	return &Handler{
		ds:     ds,
		logger: logger,
		parser: parserInstance,
	}, nil
}

// RegisterRoutes registers all assistant-related routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	parserRouter := router.PathPrefix("/parser").Subrouter()

	parserRouter.HandleFunc("/prescription", h.ParsePrescription).Methods("POST")
	parserRouter.HandleFunc("/prescription/sample", h.SaveSamplePrescription).Methods("POST")
	parserRouter.HandleFunc("/prescription/{id}", h.GetJobStatus).Methods("GET")
}
