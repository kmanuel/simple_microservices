package api

import "github.com/jinzhu/gorm"

type TaskStatusList struct {
	ID    string        `jsonapi:"primary,task_status_list"`
	Tasks []*TaskStatus `jsonapi:"attr,tasks"`
}

type TaskStatus struct {
	gorm.Model `jsonapi:"attr,model"`
	ObjID      string `jsonapi:"primary,task_status"`
	TaskID     string `jsonapi:"attr,task_id"`
	Status     string `jsonapi:"attr,status"`
	TaskType   string `jsonapi:"attr,task_type"`
}
