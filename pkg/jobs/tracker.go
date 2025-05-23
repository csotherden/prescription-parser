// Package jobs provides functionality for tracking and managing asynchronous job processing.
// It offers a simple API for creating, updating, and monitoring long-running operations
// with automatic cleanup of completed jobs.
package jobs

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// JobStatus represents the current status of a processing job.
// It describes where the job is in its lifecycle.
type JobStatus string

const (
	// JobStatusPending indicates the job has been created but processing has not started.
	JobStatusPending JobStatus = "pending"

	// JobStatusProcessing indicates the job is currently being processed.
	JobStatusProcessing JobStatus = "processing"

	// JobStatusComplete indicates the job has completed successfully.
	JobStatusComplete JobStatus = "complete"

	// JobStatusFailed indicates the job encountered an error and could not complete.
	JobStatusFailed JobStatus = "failed"
)

// Job represents a generic asynchronous job with its metadata and results.
// It includes tracking information such as timing and current status.
type Job struct {
	ID          string     `json:"id"`                     // Unique identifier for the job
	Type        string     `json:"type"`                   // Type of job being processed
	Reference   string     `json:"reference"`              // Human-readable reference or description
	Status      JobStatus  `json:"status"`                 // Current status of the job
	StartedAt   time.Time  `json:"started_at"`             // When the job was created
	CompletedAt *time.Time `json:"completed_at,omitempty"` // When the job finished (if completed)
	Error       string     `json:"error,omitempty"`        // Error message if job failed
	Result      any        `json:"result"`                 // Result data from the job (if any)
}

// Tracker manages jobs throughout their lifecycle.
// It provides thread-safe access to job information and handles job cleanup.
type Tracker struct {
	jobs  map[string]*Job // Map of job ID to job information
	mutex sync.RWMutex    // Mutex to protect concurrent access
}

// NewTracker creates a new job tracker.
// It initializes the tracker and starts the background cleanup process.
func NewTracker() *Tracker {
	tracker := &Tracker{
		jobs: make(map[string]*Job),
	}

	// Start automatic cleanup
	go tracker.startCleanup()

	return tracker
}

// startCleanup starts a background goroutine that cleans up old jobs.
// It runs periodically to remove completed jobs that are older than a defined threshold.
func (t *Tracker) startCleanup() {
	for {
		time.Sleep(5 * time.Minute)
		t.CleanupOldJobs(15 * time.Minute) // Keep jobs for 15 minutes
	}
}

// CreateJob creates a new job with the specified type and reference.
// It returns the generated job ID that can be used to check status later.
func (t *Tracker) CreateJob(jobType, reference string) string {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	jobID := uuid.New().String()
	t.jobs[jobID] = &Job{
		ID:        jobID,
		Type:      jobType,
		Reference: reference,
		Status:    JobStatusPending,
		StartedAt: time.Now(),
	}
	return jobID
}

// GetJob returns a job by ID.
// It returns the job object and a boolean indicating if the job exists.
func (t *Tracker) GetJob(jobID string) (*Job, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	job, exists := t.jobs[jobID]
	return job, exists
}

// UpdateJob updates a job's status, error message, and result.
// It returns true if the job was found and updated, false if the job doesn't exist.
func (t *Tracker) UpdateJob(jobID string, status JobStatus, err error, result any) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	job, exists := t.jobs[jobID]
	if !exists {
		return false
	}

	job.Status = status

	if status == JobStatusComplete || status == JobStatusFailed {
		now := time.Now()
		job.CompletedAt = &now
	}

	if err != nil {
		job.Error = err.Error()
	}

	job.Result = result

	return true
}

// CleanupOldJobs removes jobs older than the specified duration.
// It only removes jobs that have been completed or failed.
func (t *Tracker) CleanupOldJobs(olderThan time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	threshold := time.Now().Add(-olderThan)
	cleaned := 0

	for id, job := range t.jobs {
		if job.CompletedAt != nil && job.CompletedAt.Before(threshold) {
			delete(t.jobs, id)
			cleaned++
		}
	}

	if cleaned > 0 {
		log.Printf("Cleaned up %d old jobs", cleaned)
	}
}

// GlobalTracker is a global instance that can be used across the application.
// It provides a convenient shared job tracker accessible to all components.
var GlobalTracker = NewTracker()
