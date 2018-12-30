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
		"download_result": fmt.Sprintf("localhost:8080/images/" + t.ID),
	}
}

type ScreenShotTask struct {
	ID		string		`jsonapi:"primary,screenshot_task"`
	Url		string		`jsonapi:"attr,url"`
}

type CropTask struct {
	ID			string		`jsonapi:"primary,crop_task"`
	ImageId		string		`jsonapi:"attr,image_id"`
	Width 		int			`jsonapi:"attr,width"`
	Height 		int			`jsonapi:"attr,height"`
}
