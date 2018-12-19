package model

import (
	"encoding/json"
	"errors"
)

type Task struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"tasktype"`
	TaskParams map[string]interface{} `json:"taskParams"`
}

func (t *Task) UnmarshalJSON(data []byte) error {
	var jsonMap map[string]interface{}

	if t == nil {
		return errors.New("RawString: UnmarshalJSON on nil pointer")
	}

	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return err
	}

	t.ID = jsonMap["id"].(string)
	t.Type = jsonMap["tasktype"].(string)

	t.TaskParams = make(map[string]interface{})

	for key, val := range jsonMap {
		if key != "id" && key != "tasktype" {
			t.TaskParams[key] = val
		}
	}

	return nil
}