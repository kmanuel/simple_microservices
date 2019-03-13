package api_root

import (
	"fmt"
	"github.com/google/jsonapi"
)

type RootResource struct {
	ID			int				`jsonapi:"primary,resources"`
	Title		string			`jsonapi:"attr,title"`
	Endpoints	[]*Endpoint		`jsonapi:"relation,endpoints"`
}

type Endpoint struct {
	ID			int				`jsonapi:"primary,endpoints"`
	Path 		string			`jsonapi:"attr,path"`
}

func (r RootResource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("localhost:8080/"),
	}
}
