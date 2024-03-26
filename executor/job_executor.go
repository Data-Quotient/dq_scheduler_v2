package executor

import (
	"log"
	"net/http"
)

type JobExecutor struct{}

func NewJobExecutor() *JobExecutor {
	return &JobExecutor{}
}

func (e *JobExecutor) ExecuteJob(endpoint string) {
	resp, err := http.Post(endpoint, "application/json", nil)
	if err != nil {
		log.Printf("Error triggering job: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Job triggered successfully: %s", endpoint)
}
