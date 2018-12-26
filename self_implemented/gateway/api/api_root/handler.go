package api_root

import (
	"fmt"
	"github.com/google/jsonapi"
	"net/http"
)

const (
	headerAccept      = "Accept"
	headerContentType = "Content-Type"
)

type RootHandler struct{}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//if r.Header.Get(headerAccept) != api.MediaType {
	//	http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
	//}

	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodGet:
		methodHandler = h.getRootResource
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *RootHandler) getRootResource(w http.ResponseWriter, r *http.Request) {
	jsonapiRuntime := jsonapi.NewRuntime().Instrument("resources.list")

	resources := fixtureRootResource()
	fmt.Println("loaded resources", resources)


	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	if err := jsonapiRuntime.MarshalPayload(w, resources); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func fixtureRootResource() interface{} {
	return &Resource{
		ID:		1,
		Title: "RootResource",
		Endpoints: []*Endpoint{
			{
				ID:	1,
				Path: "/tasks",
			},
			{
				ID: 2,
				Path: "/faktory/info",
			},
			{
				ID: 3,
				Path: "/upload",
			},
			{
				ID: 4,
				Path: "/tasks/:taskId/download",
			},
		},
	}
}