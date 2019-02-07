package service

import (
	faktory "github.com/contribsys/faktory/client"
	"github.com/prometheus/common/log"
)

type FaktoryService interface {
	PublishToFaktory(taskType string, jsonTask string) error
	Info() (map[string]interface{}, error)
}

type faktoryServiceImpl struct {}

func NewFaktoryService() FaktoryService {
	return faktoryServiceImpl{}
}

//func (s faktoryServiceImpl) Publish(t *model.Task) error {
//	t.ID = uuid.New().String()
//	buf := new(bytes.Buffer)
//
//	err := jsonapi.MarshalPayload(buf, t)
//	if err != nil {
//		return err
//	}
//
//	err = s.publishToFaktory(s.taskType, buf.String())
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (faktoryServiceImpl) PublishToFaktory(taskType string, jsonTask string) error {
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

func (faktoryServiceImpl) Info() (map[string]interface{}, error) {
	client, err := faktory.Open()
	if err != nil {
		panic(err)
	}

	return client.Info()
}
