package api_root

import (
	"fmt"
	"github.com/google/jsonapi"
)

type Resource struct {
	ID			int				`api:"primary,resources"`
	Title		string			`api:"attr,title"`
	Endpoints	[]*Endpoint		`api:"relation,endpoints"`
}

type Endpoint struct {
	ID			int				`api:"primary,endpoints"`
	Path 		string			`api:"attr,path"`
}

// JSONAPILinks implements the Linkable interface for a blog
func (r Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("https://example.com/"),
	}
}

// JSONAPIRelationshipLinks implements the RelationshipLinkable interface for a blog
func (r Resource) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	//if relation == "endpoints" {
	//	return &api.Links{
	//		"related": fmt.Sprintf("https://example.com/blogs/%d/posts", 123),
	//	}
	//}
	return nil
}

// JSONAPIRelationshipMeta implements the RelationshipMetable interface for a blog
func (r Resource) JSONAPIRelationshipMeta(relation string) *jsonapi.Meta {
	if relation == "endpoints" {
		return &jsonapi.Meta{
			"detail": "endpoints offered by the gateway application",
		}
	}
	return nil
}
