package update_status

import (
	"bytes"
	"github.com/google/jsonapi"
	"github.com/prometheus/common/log"
	"net/http"
)

type TaskStatus struct {
	ObjID    string `jsonapi:"primary,task_status"`
	TaskID   string `jsonapi:"attr,task_id"`
	Status   string `jsonapi:"attr,status"`
	TaskType string `jsonapi:"attr,task_type"`
}

func NotifyAboutProcessingStart(taskId string) {
	updateTaskStatus(taskId, "processing")
}

func NotifyAboutCompletion(taskId string) {
	updateTaskStatus(taskId, "completed")
}

func updateTaskStatus(taskId string, newStatus string) error {
	taskStatus := &TaskStatus{
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
		log.Error("failed to post new TaskStatus", err)
	}
	log.Info("response=", resp)
	return err
}
