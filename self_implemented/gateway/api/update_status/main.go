package update_status

import (
	"bytes"
	"github.com/google/jsonapi"
	"github.com/prometheus/common/log"
	"net/http"
)

type TaskStatus struct {
	TaskID   string `jsonapi:"attr,task_id"`
	Status   string `jsonapi:"attr,status"`
	TaskType string `jsonapi:"attr,task_type"`
}

func NotifyAboutNewTask(taskId string, taskType string) error {
	return createNewTask(&TaskStatus{
		TaskID: taskId,
		TaskType: taskType,
		Status: "new",
	})
}

func createNewTask(t *TaskStatus) error {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, t); err != nil {
		return err
	}

	url := "http://request_service:8080/tasks/status/" + t.TaskID
	log.Info("sending update request to ", url)
	_, err := http.Post(url, jsonapi.MediaType, buf)
	if err != nil {
		log.Error("failed to post new TaskStatus", err)
	}
	return err
}

func NotifyAboutProcessingStart(taskId string) error {
	return updateTask(taskId, "processing")
}

func NotifyAboutCompletion(taskId string) error {
	return updateTask(taskId, "completed")
}

func updateTask(taskId string, newStatus string) error {
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
	_, err := http.Post(url, jsonapi.MediaType, buf)
	return err
}
