package parser

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/csotherden/prescription-parser/pkg/jobs"
	"github.com/csotherden/prescription-parser/pkg/mocks"
	"github.com/csotherden/prescription-parser/pkg/parser"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func TestParsePrescription(t *testing.T) {
	// Create a test logger
	logger := zap.NewNop()

	// Create mocks
	mockParser := mocks.NewMockParser()
	mockDatastore := mocks.NewMockDatastore()

	// Set up a job in the tracker
	createdJobID := jobs.GlobalTracker.CreateJob(parser.JobTypeParsePrescription, "Processing image: test.pdf")

	// Configure mock parser to return our job ID
	mockParser.SetParseImageResponse("test.pdf", createdJobID, nil)

	// Create test handler
	handler := NewHandler(mockParser, mockDatastore, logger)

	// Set up test router
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Create a test file for the multipart request
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.pdf")
	part.Write([]byte("test data"))
	writer.Close()

	// Create test server
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Create test request
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/parser/prescription", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse the response
	var response *jobs.Job
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify response
	if response.ID != createdJobID {
		t.Errorf("Expected job ID %s, got %s", createdJobID, response.ID)
	}

	// Verify the parser was called
	calls := mockParser.GetParseImageCalls()
	if len(calls) != 1 {
		t.Errorf("Expected 1 parser call, got %d", len(calls))
	}
}
