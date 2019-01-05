package service

import (
	"bytes"
	faktory "github.com/contribsys/faktory/client"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/model"
	"github.com/prometheus/common/log"
)

type FaktoryService interface{
	PublishTask(task model.CropTask) error
}

type FaktoryServiceImpl struct {}

func (FaktoryServiceImpl) PublishTask(task model.CropTask) error {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, task); err != nil {
		return err
	}

	log.Info("publish to faktory")
	client, err := faktory.Open()
	if err != nil {
		log.Error("failed to open connection to faktory", err)
		return err
	}
	job := faktory.NewJob("crop", buf.String())
	job.Queue = "crop"
	log.Info("publishing job", job)
	err = client.Push(job)

	return err
}
