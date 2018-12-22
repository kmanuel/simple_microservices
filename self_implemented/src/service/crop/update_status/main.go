package update_status

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type TaskStatusUpdate struct {
	Status	string	`json:"status"`
}

func NotifyAboutProcessingStart(taskId string) {
	updateTaskStatus(taskId, "processing")
}

func NotifyAboutCompletion(taskId string) {
	updateTaskStatus(taskId, "completed")
}

func updateTaskStatus(taskId string, newStatus string) {
	var nt TaskStatusUpdate
	nt.Status = newStatus
	marshal, e := json.Marshal(nt)
	if e != nil {
		panic(e)
	}
	log.Error(" sending update request for taskId", taskId)
	http.Post("http://request_service:8080/tasks/"+taskId+"/status", "application/json", bytes.NewBuffer([]byte(marshal)))
}

