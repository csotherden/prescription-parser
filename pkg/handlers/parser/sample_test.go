package parser

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/csotherden/prescription-parser/pkg/mocks"
	"github.com/csotherden/prescription-parser/pkg/models"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func TestSaveSamplePrescription(t *testing.T) {
	// Create a test logger
	logger := zap.NewNop()

	// Create mocks
	mockParser := mocks.NewMockParser()
	mockDatastore := mocks.NewMockDatastore()

	// Create a test prescription
	prescription := models.Prescription{
		Medications: []models.Medication{
			{
				DrugName: "Metformin",
				Strength: "500mg",
				Quantity: "30",
				Refills:  "3",
				SIG:      "Take 1 tablet by mouth daily",
			},
		},
		Patient: models.Patient{
			FirstName: "John",
			LastName:  "Doe",
			Dob:       "1990-01-01",
		},
		Prescriber: models.Prescriber{
			Name:      "Jane Smith",
			Npi:       "1234567890",
			Specialty: "MD",
		},
	}

	// Convert prescription to JSON
	prescriptionJSON, err := json.Marshal(prescription)
	if err != nil {
		t.Fatalf("Failed to marshal prescription: %v", err)
	}

	// Set expected values
	testImageID := "test-image-id"
	testEmbedding := []float32{0.1, 0.2, 0.3, 0.4}

	// Configure mock parser responses
	mockParser.SetUploadImageResponse("test.pdf", testImageID, nil)
	mockParser.SetEmbedding(prescription.Medications[0].DrugName, testEmbedding, nil)

	// Create test handler
	handler := NewHandler(mockParser, mockDatastore, logger)

	// Set up test router
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Create a test server
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Create a multipart form with test data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add prescription JSON
	err = writer.WriteField("json", string(prescriptionJSON))
	if err != nil {
		t.Fatalf("Failed to add JSON field: %v", err)
	}

	// Add test file
	part, err := writer.CreateFormFile("image", "test.pdf")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write([]byte("test pdf content"))
	writer.Close()

	// Create and send the request
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/parser/prescription/sample", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	// Check that the parser methods were called
	uploadCalls := mockParser.GetUploadImageCalls()
	if len(uploadCalls) != 1 {
		t.Errorf("Expected 1 UploadImage call, got %d", len(uploadCalls))
	}

	embeddingCalls := mockParser.GetEmbeddingCalls()
	if len(embeddingCalls) != 1 {
		t.Errorf("Expected 1 GetEmbedding call, got %d", len(embeddingCalls))
	}

	// Verify datastore interaction
	saveCalls := mockDatastore.GetSaveSamplePrescriptionCalls()
	if len(saveCalls) != 1 {
		t.Errorf("Expected 1 SaveSamplePrescription call, got %d", len(saveCalls))
	}

	// Verify parameters passed to datastore
	if len(saveCalls) > 0 {
		call := saveCalls[0]

		// Check image ID
		if call.ImageID != testImageID {
			t.Errorf("Expected image ID %s, got %s", testImageID, call.ImageID)
		}

		// Check MIME type
		expectedMimeType := "application/pdf"
		if call.MimeType != expectedMimeType {
			t.Errorf("Expected MIME type %s, got %s", expectedMimeType, call.MimeType)
		}

		// Check prescription (using the first medication's drug name as identifier)
		if len(call.Prescription.Medications) == 0 ||
			call.Prescription.Medications[0].DrugName != prescription.Medications[0].DrugName {
			t.Errorf("Prescription data mismatch")
		}

		// Check embedding
		if len(call.Embedding) != len(testEmbedding) {
			t.Errorf("Expected embedding length %d, got %d", len(testEmbedding), len(call.Embedding))
		}
	}
}
