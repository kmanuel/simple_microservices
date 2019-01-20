package service

import (
	"bytes"
	faktory "github.com/contribsys/faktory/client"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/simple_microservices/self_implemented/service/portrait/model"
	"github.com/prometheus/common/log"
)

type FaktoryPublishService interface {
	Publish(t *model.Task) error
}

type faktoryPublishService struct {
	taskType string
}

func NewFaktoryPublishService(taskType string) FaktoryPublishService {
	return faktoryPublishService{taskType: taskType}
}

func (s faktoryPublishService) Publish(t *model.Task) error {
	t.ID = uuid.New().String()
	buf := new(bytes.Buffer)

	err := jsonapi.MarshalPayload(buf, t)
	if err != nil {
		return err
	}

	err = s.publishToFaktory(s.taskType, buf.String())
	if err != nil {
		return err
	}

	return nil
}

func (faktoryPublishService) publishToFaktory(taskType string, jsonTask string) error {
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
