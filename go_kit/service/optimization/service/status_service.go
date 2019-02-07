package service

import (
	"bytes"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/google/jsonapi"
	"github.com/prometheus/common/log"
	"net/http"
)

type taskStatus struct {
	TaskID   string `jsonapi:"attr,task_id"`
	Status   string `jsonapi:"attr,status"`
	TaskType string `jsonapi:"attr,task_type"`
}

type StatusClient interface {
	NotifyAboutProcessing(taskId string) error
	NotifyAboutCompletion(taskId string) error
	NotifyAboutFailure(taskId string) error
}

type statusClientImpl struct {}

func NewStatusClient() StatusClient {
	return statusClientImpl{}
}

func (statusClientImpl) NotifyAboutProcessing(taskId string) error {
	return updateTaskStatus(taskId, "processing")
}

func (statusClientImpl) NotifyAboutCompletion(taskId string) error {
	return updateTaskStatus(taskId, "completed")
}

func (statusClientImpl) NotifyAboutFailure(taskId string) error {
	return updateTaskStatus(taskId, "failed")
}

func updateTaskStatus(taskId string, newStatus string) error {
	return hystrix.Do("update_task_status", func() error {
		taskStatus := &taskStatus{
			TaskID:   taskId,
			Status:   newStatus,
		}
		buf := new(bytes.Buffer)
		if err := jsonapi.MarshalPayload(buf, taskStatus); err != nil {
			return err
		}

		url := "http://request_service:8080/tasks/status/" + taskId
		log.Info("sending update request to ", url)
		resp, err := http.Post(url, jsonapi.MediaType, buf)
		if err != nil {
			log.Error("failed to post new taskStatus", err)
		}
		log.Info("response=", resp)
		return err
	}, nil)

}