package api_root

import (
	"fmt"
	"github.com/google/jsonapi"
	"net/http"
)

type RootHandler struct{}

func (h *RootHandler) GetRootResource(w http.ResponseWriter, r *http.Request) {
	resources := fixtureRootResource(r.Host)
	fmt.Println("loaded resources", resources)


	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	if err := jsonapi.MarshalPayload(w, resources); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

