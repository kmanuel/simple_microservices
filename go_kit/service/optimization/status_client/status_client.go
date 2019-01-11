package status_client

import (
	"bytes"
	"github.com/google/jsonapi"
	"github.com/prometheus/common/log"
	"net/http"
	"os"
)

type taskStatus struct {
	TaskID   string `jsonapi:"attr,task_id"`
	Status   string `jsonapi:"attr,status"`
	TaskType string `jsonapi:"attr,task_type"`
}

type StatusClient interface {
	NotifyAboutNew(taskId string) error
	NotifyAboutProcessing(taskId string) error
	NotifyAboutCompletion(taskId string) error
	NotifyAboutFailure(taskId string) error
}

type statusClientImpl struct {
	TaskType string
}

func NewStatusClient(taskType string) StatusClient {
	return statusClientImpl{TaskType: taskType}
}

func (s statusClientImpl) NotifyAboutNew(taskId string) error {
	return s.updateTaskStatus(taskId, "new")
}

func (s statusClientImpl) NotifyAboutProcessing(taskId string) error {
	return s.updateTaskStatus(taskId, "processing")
}

func (s statusClientImpl) NotifyAboutCompletion(taskId string) error {
	return s.updateTaskStatus(taskId, "completed")
}

func (s statusClientImpl) NotifyAboutFailure(taskId string) error {
	return s.updateTaskStatus(taskId, "failed")
}

func (s statusClientImpl) updateTaskStatus(taskId string, newStatus string) error {
	taskStatus := &taskStatus{
		TaskID:   taskId,
		Status:   newStatus,
		TaskType: s.TaskType,
	}
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, taskStatus); err != nil {
		return err
	}

	requestServiceHost := os.Getenv("REQUEST_SERVICE_HOST") + ":" + os.Getenv("REQUEST_SERVICE_PORT")
	url := requestServiceHost + "/tasks/status/" + taskId
	log.Info("sending update request to ", url)
	resp, err := http.Post(url, jsonapi.MediaType, buf)
	if err != nil {
		log.Error("failed to post new taskStatus", err)
	}
	log.Info("response=", resp)
	return err
}
