package api_image

import (
	"fmt"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/api/api_task"
)

type Image struct {
	ID			string				`jsonapi:"primary,images"`
	Tasks		[]api_task.Task		`jsonapi:"relation,tasks"`
}

func (img Image) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("localhost:8080/images/" + img.ID),
	}
}
