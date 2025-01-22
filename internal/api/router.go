// internal/api/router.go
package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/submit", SubmitJob)
		api.GET("/status", GetJobStatus)
	}

	return r
}
