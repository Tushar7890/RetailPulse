// internal/api/handlers.go
package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Tushar7890/RetailPulse/internal/jobs"
	"github.com/Tushar7890/RetailPulse/internal/models"
)

var jobManager = jobs.NewManager()

func SubmitJob(c *gin.Context) {
	var payload models.JobRequest
	if err := c.ShouldBindJSON(&payload); err != nil || len(payload.Visits) != payload.Count {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	jobID := uuid.New().String()
	job := jobManager.CreateJob(jobID, payload)

	go jobManager.ProcessJob(job) // Run processing in background

	c.JSON(http.StatusCreated, gin.H{"job_id": jobID})
}

func GetJobStatus(c *gin.Context) {
	jobID := c.Query("jobid")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID required"})
		return
	}

	status, err := jobManager.GetJobStatus(jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
