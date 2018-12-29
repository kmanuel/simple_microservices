package api_task

import (
	"fmt"
	"github.com/google/jsonapi"
)

type Task struct {
	ID         string				  `jsonapi:"primary,tasks"`
	Type       string                 `jsonapi:"attr,task_type"`
	TaskParams map[string]interface{} `jsonapi:"attr,task_params"`
}

func (t *Task) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("localhost:8080/tasks/" + t.ID),
	}
}
