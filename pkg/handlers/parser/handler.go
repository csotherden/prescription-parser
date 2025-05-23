package parser

import (
	"github.com/csotherden/prescription-parser/pkg/datastore"
	parserPkg "github.com/csotherden/prescription-parser/pkg/parser"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Handler handles parser-related API requests
type Handler struct {
	ds     *datastore.Datastore
	logger *zap.Logger
	parser parserPkg.Parser
}

// NewHandler creates a new parser handler instance
func NewHandler(parser parserPkg.Parser, ds *datastore.Datastore, logger *zap.Logger) *Handler {
	return &Handler{
		ds:     ds,
		logger: logger,
		parser: parser,
	}
}

// RegisterRoutes registers all parser-related routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	parserRouter := router.PathPrefix("/parser").Subrouter()

	parserRouter.HandleFunc("/prescription", h.ParsePrescription).Methods("POST")
	parserRouter.HandleFunc("/prescription/sample", h.SaveSamplePrescription).Methods("POST")
	parserRouter.HandleFunc("/prescription/{id}", h.GetJobStatus).Methods("GET")
}
