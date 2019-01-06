package service

import (
	"bytes"
	"fmt"
	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/model"
	"github.com/prometheus/common/log"
)

type FaktoryPublishService interface {
	PublishTask(task *model.CropTask) error
}

type FaktoryListenService interface{
	Handle(queue string, fn worker.Perform)
}

type FaktoryService interface {
	PublishTask(task *model.CropTask) error
	Handle(queue string, fn worker.Perform)
}

func NewFaktoryService(cropService CropService) FaktoryService {
	return faktoryServiceImpl{cropService}
}

type faktoryServiceImpl struct {
	cropService CropService
}

func (faktoryServiceImpl) PublishTask(task *model.CropTask) error {
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

func (fs faktoryServiceImpl) Handle(queue string, fn worker.Perform) {
	fmt.Println("starting faktory")
	mgr := worker.NewManager()
	mgr.Register(queue, fn)
	mgr.Queues = []string{queue}
	var quit bool
	mgr.On(worker.Shutdown, func() {
		quit = true
	})
	go mgr.Run()
}

