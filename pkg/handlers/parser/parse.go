package parser

import (
	"fmt"
	"net/http"

	"github.com/csotherden/prescription-parser/pkg/handlerutils"
	"github.com/csotherden/prescription-parser/pkg/jobs"
	"go.uber.org/zap"
)

const (
	maxUploadSize = 10 << 20 // 10 MB
)

// ParsePrescription handles the request to parse a new prescription image
func (h *Handler) ParsePrescription(w http.ResponseWriter, r *http.Request) {
	// Validate file size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		handlerutils.RespondWithError(w, h.logger, http.StatusBadRequest, "File exceeds max upload size", fmt.Errorf("file exceeds max upload size: %w", err))
		return
	}

	// Get the image file from the request
	file, header, err := r.FormFile("image")
	if err != nil {
		handlerutils.RespondWithError(w, h.logger, http.StatusBadRequest, fmt.Sprintf("Failed to read image file: %s", err.Error()), fmt.Errorf("failed to read image file: %w", err))
		return
	}
	defer file.Close()

	jobID, err := h.parser.ParseImage(r.Context(), header.Filename, file)
	if err != nil {
		h.logger.Error("failed to parse image", zap.Error(err))
		handlerutils.RespondWithError(w, h.logger, http.StatusInternalServerError, "failed to process image", err)
	}

	job, exists := jobs.GlobalTracker.GetJob(jobID)
	if !exists || job == nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	handlerutils.RespondWithJSON(w, h.logger, http.StatusOK, job)
}
