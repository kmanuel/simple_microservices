package api_task

import (
	faktory "github.com/contribsys/faktory/client"
	"github.com/prometheus/common/log"
)

func publishToFactory(t *Task) error {
	log.Info("publish to faktory", t)
	client, err := faktory.Open()
	if err != nil {
		log.Error("failed to open connection to faktory", err)
		return err
	}
	job := faktory.NewJob(t.Type, &t.TaskParams)
	job.Queue = t.Type
	t.TaskParams["id"] = t.ID
	job.Custom = t.TaskParams
	log.Info("publishing job", job)
	err = client.Push(job)

	return err
}

func publishCropToFactory(taskType string, jsonTask string) error {
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
