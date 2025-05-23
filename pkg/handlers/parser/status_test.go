package parser

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/csotherden/prescription-parser/pkg/jobs"
	"github.com/csotherden/prescription-parser/pkg/mocks"
	"github.com/csotherden/prescription-parser/pkg/parser"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func TestGetJobStatus(t *testing.T) {
	// Create a test logger
	logger := zap.NewNop()

	// Create mocks
	mockParser := mocks.NewMockParser()
	mockDatastore := mocks.NewMockDatastore()

	// Create a handler
	handler := NewHandler(mockParser, mockDatastore, logger)

	// Set up test router
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Create test jobs
	completedJobID := jobs.GlobalTracker.CreateJob(parser.JobTypeParsePrescription, "Processing image: completed.pdf")
	pendingJobID := jobs.GlobalTracker.CreateJob(parser.JobTypeParsePrescription, "Processing image: pending.pdf")

	// Update job status to complete
	jobs.GlobalTracker.UpdateJob(completedJobID, jobs.JobStatusComplete, nil, nil)

	// Create a test server
	ts := httptest.NewServer(router)
	defer ts.Close()

	tests := []struct {
		name           string
		jobID          string
		expectedStatus int
		checkBody      bool
	}{
		{
			name:           "Job found - completed",
			jobID:          completedJobID,
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "Job found - pending",
			jobID:          pendingJobID,
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "Job not found",
			jobID:          "non-existent-id",
			expectedStatus: http.StatusNotFound,
			checkBody:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			url := ts.URL + "/parser/prescription/" + tt.jobID
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Send request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// For successful requests, verify the job details
			if tt.checkBody && resp.StatusCode == http.StatusOK {
				var job *jobs.Job
				err = json.NewDecoder(resp.Body).Decode(&job)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				// Check that the returned job has the correct ID
				if job.ID != tt.jobID {
					t.Errorf("Expected job ID %s, got %s", tt.jobID, job.ID)
				}

				// Check the job status
				if tt.jobID == completedJobID && job.Status != jobs.JobStatusComplete {
					t.Errorf("Expected job status %s, got %s", jobs.JobStatusComplete, job.Status)
				} else if tt.jobID == pendingJobID && job.Status != jobs.JobStatusPending {
					t.Errorf("Expected job status %s, got %s", jobs.JobStatusPending, job.Status)
				}
			}
		})
	}
}
