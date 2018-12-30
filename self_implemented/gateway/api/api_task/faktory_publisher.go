package api_task

import (
	faktory "github.com/contribsys/faktory/client"
	"github.com/prometheus/common/log"
)

func publishToFaktory(taskType string, jsonTask string) error {
	log.Info("publish to faktory")
	client, err := faktory.Open()
	if err != nil {
		log.Error("failed to open connection to faktory", err)
		return err
	}
	job := faktory.NewJob(taskType, jsonTask)
	job.Queue = taskType
	log.Info("publishing job", job)
	err = client.Push(job)

	return err
}
