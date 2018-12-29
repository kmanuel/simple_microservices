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

// JSONAPILinks implements the Linkable interface for a blog
func (r RootResource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("localhost:8080/"),
	}
}
//
//// JSONAPIRelationshipLinks implements the RelationshipLinkable interface for a blog
//func (r RootResource) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
//	if relation == "endpoints" {
//		return &jsonapi.Links{
//			"related": fmt.Sprintf("https://example.com/blogs/%d/posts", 123),
//		}
//	}
//	return nil
//}
//
//// JSONAPIRelationshipMeta implements the RelationshipMetable interface for a blog
//func (r RootResource) JSONAPIRelationshipMeta(relation string) *jsonapi.Meta {
//	if relation == "endpoints" {
//		return &jsonapi.Meta{
//			"detail": "endpoints offered by the gateway application",
//		}
//	}
//	return nil
//}
