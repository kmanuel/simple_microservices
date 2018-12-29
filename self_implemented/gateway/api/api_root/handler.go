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

	resources := fixtureRootResource(r.Host)
	fmt.Println("loaded resources", resources)


	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	if err := jsonapiRuntime.MarshalPayload(w, resources); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

