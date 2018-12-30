package api

import "github.com/jinzhu/gorm"

type NewTaskDTO struct {
	TaskId string `json:"id"`
}

type TaskStatusUpdateDTO struct {
	Status string `json:"status"`
}

type TaskStatusList struct {
	ID 		string			`jsonapi:"primary,task_status_list"`
	Tasks 	[]*TaskStatus	`jsonapi:"relation,tasks"`
}

type TaskStatus struct {
	gorm.Model		`jsonapi:"attr,model"`
	TaskId string 	`jsonapi:"primary,task_status"`
	Status string 	`jsonapi:"attr,status"`
}
