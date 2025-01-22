// internal/jobs/processor.go
package jobs

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/Tushar7890/RetailPulse/internal/models"
)

type Manager struct {
	jobs map[string]*models.Job
	mu   sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		jobs: make(map[string]*models.Job),
	}
}

func (m *Manager) CreateJob(jobID string, payload models.JobRequest) *models.Job {
	m.mu.Lock()
	defer m.mu.Unlock()

	job := &models.Job{
		ID:       jobID,
		Visits:   payload.Visits,
		Status:   "ongoing",
		Errors:   []models.JobError{},
		Results:  []models.ImageResult{},
		CreatedAt: time.Now(),
	}
	m.jobs[jobID] = job
	return job
}

func (m *Manager) ProcessJob(job *models.Job) {
	for _, visit := range job.Visits {
		for _, imageURL := range visit.ImageURL {
			resp, err := http.Get(imageURL)
			if err != nil {
				job.Errors = append(job.Errors, models.JobError{StoreID: visit.StoreID, Error: "Failed to download image"})
				continue
			}
			img, err := imaging.Decode(resp.Body)
			resp.Body.Close()
			if err != nil {
				job.Errors = append(job.Errors, models.JobError{StoreID: visit.StoreID, Error: "Invalid image format"})
				continue
			}

			bounds := img.Bounds()
			perimeter := 2 * (bounds.Dx() + bounds.Dy())
			time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)

			job.Results = append(job.Results, models.ImageResult{
				StoreID:   visit.StoreID,
				ImageURL:  imageURL,
				Perimeter: perimeter,
			})
		}
	}
	if len(job.Errors) > 0 {
		job.Status = "failed"
	} else {
		job.Status = "completed"
	}
}

func (m *Manager) GetJobStatus(jobID string) (map[string]interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, exists := m.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job not found")
	}

	return map[string]interface{}{
		"status": job.Status,
		"job_id": job.ID,
		"results": job.Results,
		"errors": job.Errors,
	}, nil
}
