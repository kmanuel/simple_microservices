package api_faktory

import (
	"encoding/json"
	faktory "github.com/contribsys/faktory/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"net/http"
)

type FaktoryHandler struct{
	RequestCounter *prometheus.CounterVec
}

func (h *FaktoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodGet:
		methodHandler = h.getFaktoryInfo
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *FaktoryHandler) getFaktoryInfo(w http.ResponseWriter, r *http.Request) {
	client, err := faktory.Open()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	info, err := client.Info()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	err = json.NewEncoder(w).Encode(info)
	if err != nil {
		log.Error("error writing response")
	}
}
