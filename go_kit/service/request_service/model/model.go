package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type TaskStatusList struct {
	ID    string        `jsonapi:"primary,task_status_list"`
	Tasks []*TaskStatus `jsonapi:"attr,tasks"`
}

type TaskStatus struct {
	gorm.Model `jsonapi:"attr,model"`
	TaskID     string `jsonapi:"attr,task_id"`
	Status     string `jsonapi:"attr,status"`
	TaskType   string `jsonapi:"attr,task_type"`
}

func (ts TaskStatus) String() string {
	return fmt.Sprintf("{TaskStatus[TaskID=%s][Status=%s][TaskType=%s]}", ts.TaskID, ts.Status, ts.TaskType)
}
