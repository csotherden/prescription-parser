package mocks

import (
	"context"
	"fmt"
	"sync"

	"github.com/csotherden/prescription-parser/pkg/models"
)

// MockDatastore implements the datastore.Datastore interface for testing
type MockDatastore struct {
	mu                          sync.Mutex
	getSamplesCalls             []getSamplesCall
	saveSamplePrescriptionCalls []saveSamplePrescriptionCall
	samples                     map[string][]models.SamplePrescription
	samplesErr                  map[string]error
	savedSamples                map[string]models.Prescription
	saveSampleErr               map[string]error
}

type getSamplesCall struct {
	Ctx       context.Context
	Embedding []float32
}

type saveSamplePrescriptionCall struct {
	Ctx          context.Context
	MimeType     string
	ImageID      string
	Prescription models.Prescription
	Embedding    []float32
}

// NewMockDatastore creates a new mock datastore
func NewMockDatastore() *MockDatastore {
	return &MockDatastore{
		samples:       make(map[string][]models.SamplePrescription),
		samplesErr:    make(map[string]error),
		savedSamples:  make(map[string]models.Prescription),
		saveSampleErr: make(map[string]error),
	}
}

// GetSamples mocks the GetSamples method
func (m *MockDatastore) GetSamples(ctx context.Context, embedding []float32) ([]models.SamplePrescription, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create a simple key based on the first few elements of the embedding
	key := createEmbeddingKey(embedding)

	m.getSamplesCalls = append(m.getSamplesCalls, getSamplesCall{
		Ctx:       ctx,
		Embedding: embedding,
	})

	if err, ok := m.samplesErr[key]; ok && err != nil {
		return nil, err
	}

	samples, ok := m.samples[key]
	if !ok {
		// Return empty slice if no samples are set
		return []models.SamplePrescription{}, nil
	}

	return samples, nil
}

// SaveSamplePrescription mocks the SaveSamplePrescription method
func (m *MockDatastore) SaveSamplePrescription(ctx context.Context, mimeType, imageID string,
	prescription models.Prescription, embedding []float32) error {

	m.mu.Lock()
	defer m.mu.Unlock()

	key := imageID

	m.saveSamplePrescriptionCalls = append(m.saveSamplePrescriptionCalls, saveSamplePrescriptionCall{
		Ctx:          ctx,
		MimeType:     mimeType,
		ImageID:      imageID,
		Prescription: prescription,
		Embedding:    embedding,
	})

	if err, ok := m.saveSampleErr[key]; ok && err != nil {
		return err
	}

	m.savedSamples[key] = prescription

	return nil
}

// SetSamplePrescriptions configures the mock to return specific sample prescriptions for a given embedding
func (m *MockDatastore) SetSamplePrescriptions(embedding []float32, samples []models.SamplePrescription, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := createEmbeddingKey(embedding)
	m.samples[key] = samples

	if err != nil {
		m.samplesErr[key] = err
	}
}

// SetSaveError configures the mock to return a specific error when saving a prescription
func (m *MockDatastore) SetSaveError(imageID string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.saveSampleErr[imageID] = err
}

// GetSavedPrescription retrieves a saved prescription by imageID
func (m *MockDatastore) GetSavedPrescription(imageID string) (models.Prescription, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	prescription, ok := m.savedSamples[imageID]
	return prescription, ok
}

// GetSamplesCalls returns the recorded GetSamples calls
func (m *MockDatastore) GetSamplesCalls() []getSamplesCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.getSamplesCalls
}

// GetSaveSamplePrescriptionCalls returns the recorded SaveSamplePrescription calls
func (m *MockDatastore) GetSaveSamplePrescriptionCalls() []saveSamplePrescriptionCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.saveSamplePrescriptionCalls
}

// Helper function to create a simple key from an embedding
func createEmbeddingKey(embedding []float32) string {
	if len(embedding) == 0 {
		return "empty_embedding"
	}

	// Use first element as a simple key for test matching
	return fmt.Sprintf("%.4f", embedding[0])
}
