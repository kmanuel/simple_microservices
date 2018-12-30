package api_image

import (
	"fmt"
	"github.com/google/jsonapi"
)

type Image struct {
	ID			string				`jsonapi:"primary,images"`
}

func (img Image) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("localhost:8080/images/" + img.ID),
	}
}
