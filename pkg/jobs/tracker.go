package jobs

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// JobStatus represents the current status of a processing job
type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusComplete   JobStatus = "complete"
	JobStatusFailed     JobStatus = "failed"
)

// Job represents a generic asynchronous job
type Job struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Reference   string     `json:"reference"`
	Status      JobStatus  `json:"status"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Error       string     `json:"error,omitempty"`
	Result      any        `json:"result"`
}

// Tracker manages jobs
type Tracker struct {
	jobs  map[string]*Job
	mutex sync.RWMutex
}

// NewTracker creates a new job tracker
func NewTracker() *Tracker {
	tracker := &Tracker{
		jobs: make(map[string]*Job),
	}

	// Start automatic cleanup
	go tracker.startCleanup()

	return tracker
}

// startCleanup starts a background goroutine that cleans up old jobs
func (t *Tracker) startCleanup() {
	for {
		time.Sleep(5 * time.Minute)
		t.CleanupOldJobs(15 * time.Minute) // Keep jobs for 15 minutes
	}
}

// CreateJob creates a new job
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

// GetJob returns a job by ID
func (t *Tracker) GetJob(jobID string) (*Job, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	job, exists := t.jobs[jobID]
	return job, exists
}

// UpdateJob updates a job's status
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

// CleanupOldJobs removes jobs older than the specified duration
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

// GlobalTracker is a global instance that can be used across the application
var GlobalTracker = NewTracker()
