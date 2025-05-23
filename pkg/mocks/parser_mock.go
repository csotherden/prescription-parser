package mocks

import (
	"context"
	"io"
	"sync"

	"github.com/csotherden/prescription-parser/pkg/models"
)

// MockParser implements the parser.Parser interface for testing
type MockParser struct {
	mu                 sync.Mutex
	parseImageCalls    []parseImageCall
	getEmbeddingCalls  []getEmbeddingCall
	uploadImageCalls   []uploadImageCall
	parseImageResponse map[string]string
	parseImageErr      map[string]error
	embeddings         map[string][]float32
	embeddingErr       map[string]error
	uploadImageIDs     map[string]string
	uploadImageErr     map[string]error
}

type parseImageCall struct {
	ctx      context.Context
	fileName string
}

type getEmbeddingCall struct {
	ctx          context.Context
	prescription models.Prescription
}

type uploadImageCall struct {
	ctx      context.Context
	fileName string
}

// NewMockParser creates a new mock parser
func NewMockParser() *MockParser {
	return &MockParser{
		parseImageResponse: make(map[string]string),
		parseImageErr:      make(map[string]error),
		embeddings:         make(map[string][]float32),
		embeddingErr:       make(map[string]error),
		uploadImageIDs:     make(map[string]string),
		uploadImageErr:     make(map[string]error),
	}
}

// ParseImage mocks the ParseImage method
func (m *MockParser) ParseImage(ctx context.Context, fileName string, file io.Reader) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.parseImageCalls = append(m.parseImageCalls, parseImageCall{
		ctx:      ctx,
		fileName: fileName,
	})

	if err, ok := m.parseImageErr[fileName]; ok && err != nil {
		return "", err
	}

	jobID, ok := m.parseImageResponse[fileName]
	if !ok {
		jobID = "mock-job-id-" + fileName
		m.parseImageResponse[fileName] = jobID
	}

	return jobID, nil
}

// GetEmbedding mocks the GetEmbedding method
func (m *MockParser) GetEmbedding(ctx context.Context, prescription models.Prescription) ([]float32, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.getEmbeddingCalls = append(m.getEmbeddingCalls, getEmbeddingCall{
		ctx:          ctx,
		prescription: prescription,
	})

	// Use medication name as key for test simplicity
	key := ""
	if len(prescription.Medications) > 0 {
		key = prescription.Medications[0].DrugName
	}

	if err, ok := m.embeddingErr[key]; ok && err != nil {
		return nil, err
	}

	embedding, ok := m.embeddings[key]
	if !ok {
		// Return a default embedding if none is set
		return []float32{0.1, 0.2, 0.3}, nil
	}

	return embedding, nil
}

// UploadImage mocks the UploadImage method
func (m *MockParser) UploadImage(ctx context.Context, fileName string, file io.Reader) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.uploadImageCalls = append(m.uploadImageCalls, uploadImageCall{
		ctx:      ctx,
		fileName: fileName,
	})

	if err, ok := m.uploadImageErr[fileName]; ok && err != nil {
		return "", err
	}

	imageID, ok := m.uploadImageIDs[fileName]
	if !ok {
		imageID = "mock-image-id-" + fileName
		m.uploadImageIDs[fileName] = imageID
	}

	return imageID, nil
}

// SetParseImageResponse sets the response for a particular file name
func (m *MockParser) SetParseImageResponse(fileName, jobID string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.parseImageResponse[fileName] = jobID
	if err != nil {
		m.parseImageErr[fileName] = err
	}
}

// SetEmbedding sets the embedding for a particular prescription
func (m *MockParser) SetEmbedding(key string, embedding []float32, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.embeddings[key] = embedding
	if err != nil {
		m.embeddingErr[key] = err
	}
}

// SetUploadImageResponse sets the response for a particular file name
func (m *MockParser) SetUploadImageResponse(fileName, imageID string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.uploadImageIDs[fileName] = imageID
	if err != nil {
		m.uploadImageErr[fileName] = err
	}
}

// GetParseImageCalls returns the recorded ParseImage calls
func (m *MockParser) GetParseImageCalls() []parseImageCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.parseImageCalls
}

// GetEmbeddingCalls returns the recorded GetEmbedding calls
func (m *MockParser) GetEmbeddingCalls() []getEmbeddingCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.getEmbeddingCalls
}

// GetUploadImageCalls returns the recorded UploadImage calls
func (m *MockParser) GetUploadImageCalls() []uploadImageCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.uploadImageCalls
}
