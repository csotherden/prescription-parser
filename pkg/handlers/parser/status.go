package parser

import (
	"github.com/csotherden/prescription-parser/pkg/handlerutils"
	"github.com/csotherden/prescription-parser/pkg/jobs"
	"github.com/gorilla/mux"
	"net/http"
)

// GetJobStatus handles the request to get the status of a background job
func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	// Get job ID from URL
	vars := mux.Vars(r)
	jobID := vars["id"]
	if jobID == "" {
		handlerutils.RespondWithError(w, h.logger, http.StatusBadRequest, "Job ID is required", nil)
		return
	}

	// Get job from tracker
	job, exists := jobs.GlobalTracker.GetJob(jobID)
	if !exists || job == nil {
		handlerutils.RespondWithError(w, h.logger, http.StatusNotFound, "Job not found", nil)
		return
	}

	// Return job status
	handlerutils.RespondWithJSON(w, h.logger, http.StatusOK, job)
}
