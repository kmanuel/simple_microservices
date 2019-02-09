package service

import (
	faktory "github.com/contribsys/faktory/client"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/model"
	"github.com/prometheus/common/log"
)

type FaktoryService interface {
	PublishToFaktory(taskType string, jsonTask string) error
	Info() (*model.FaktoryInfo, error)
}

type faktoryServiceImpl struct {}

func NewFaktoryService() FaktoryService {
	return faktoryServiceImpl{}
}

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

func (faktoryServiceImpl) Info() (*model.FaktoryInfo, error) {
	client, err := faktory.Open()
	if err != nil {
		panic(err)
	}

	info, err := client.Info()

	return toFaktoryInfo(info), nil
}

func toFaktoryInfo(info map[string]interface{}) *model.FaktoryInfo {
	faktoryPart := info["faktory"].(map[string]interface{})

	queues := make(map[string]float64)

	for k, v := range faktoryPart["queues"].(map[string]interface{}) {
		queues[k] = v.(float64)
	}

	return &model.FaktoryInfo{
		TotalProcessed: faktoryPart["total_processed"].(float64),
		TotalQueues:    faktoryPart["total_queues"].(float64),
		TotalEnqueued:  faktoryPart["total_enqueued"].(float64),
		TotalFailures:  faktoryPart["total_failures"].(float64),
		Queues:         queues,
	}
}
