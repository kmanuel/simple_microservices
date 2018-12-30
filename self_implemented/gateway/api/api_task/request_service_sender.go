package api_task

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func sendToRequestService(taskId string) error {
	var nt NewTaskType
	nt.Id = taskId
	marshal, e := json.Marshal(nt)
	if e != nil {
		panic(e)
	}
	_, err := http.Post("http://request_service:8080/tasks", "application/json", bytes.NewBuffer([]byte(marshal)))
	return err
}
