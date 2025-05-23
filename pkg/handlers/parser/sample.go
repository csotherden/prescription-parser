package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/csotherden/prescription-parser/pkg/handlerutils"
	"github.com/csotherden/prescription-parser/pkg/models"
)

// SaveSamplePrescription handles the request to save a validated sample prescription
func (h *Handler) SaveSamplePrescription(w http.ResponseWriter, r *http.Request) {
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

	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	var contentType string
	switch fileExt {
	case ".pdf":
		contentType = "application/pdf"
	default:
		handlerutils.RespondWithError(w, h.logger, http.StatusBadRequest, fmt.Sprintf("unsupported file type. file must be PDF not %s", fileExt), fmt.Errorf("unsupported file type. file must be PDF not %s", fileExt))
		return
	}

	rxJson := r.FormValue("json")
	if rxJson == "" {
		handlerutils.RespondWithError(w, h.logger, http.StatusBadRequest, "Prescription JSON is required", fmt.Errorf("json is required"))
		return
	}

	var rx models.Prescription
	if err := json.Unmarshal([]byte(rxJson), &rx); err != nil {
		handlerutils.RespondWithError(w, h.logger, http.StatusBadRequest, "Invalid prescription JSON", fmt.Errorf("invalid prescription JSON: %w", err))
		return
	}

	h.logger.Info("saving sample prescription image", zap.String("file_name", header.Filename))

	imageID, err := h.parser.UploadImage(r.Context(), header.Filename, file)
	if err != nil {
		h.logger.Error("failed to upload sample image", zap.Error(err))
		handlerutils.RespondWithError(w, h.logger, http.StatusInternalServerError, "Failed to upload image", fmt.Errorf("failed to upload image: %w", err))
		return
	}

	h.logger.Info("uploaded sample prescription image", zap.String("file_name", header.Filename), zap.String("image_id", imageID))

	embedding, err := h.parser.GetEmbedding(r.Context(), rx)
	if err != nil {
		h.logger.Error("failed to generate embedding", zap.Error(err))
		handlerutils.RespondWithError(w, h.logger, http.StatusInternalServerError, "Failed to generate embedding", fmt.Errorf("failed to generate embedding: %w", err))
		return
	}

	h.logger.Info("generated embedding", zap.String("file_name", header.Filename), zap.String("image_id", imageID))

	err = h.ds.SaveSamplePrescription(r.Context(), contentType, imageID, rx, embedding)
	if err != nil {
		h.logger.Error("failed to save sample prescription", zap.Error(err))
		handlerutils.RespondWithError(w, h.logger, http.StatusInternalServerError, "Failed to save sample prescription", fmt.Errorf("failed to save sample prescription: %w", err))
		return
	}

	h.logger.Info("successfully saved sample prescription", zap.String("file_name", header.Filename), zap.String("image_id", imageID))

	handlerutils.RespondWithNoContent(w)
}
