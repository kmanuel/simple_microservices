package api_image

import (
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kmanuel/minioconnector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"io"
	"net/http"
)

type ImageHandler struct{
	RequestCounter *prometheus.CounterVec
}

func (h *ImageHandler) ServeDownload(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodGet:
		methodHandler = h.downloadImage
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *ImageHandler) downloadImage(w http.ResponseWriter, r *http.Request) {
	log.Info("download api_image request")
	imageId := mux.Vars(r)["id"]

	log.Info("download request for imageId=", imageId)
	object, err := minioconnector.GetObject(imageId)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	_, err = io.Copy(w, object)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func (h *ImageHandler) ServeUploadHTTP(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodPost:
		methodHandler = h.uploadImage
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *ImageHandler) uploadImage(w http.ResponseWriter, r *http.Request) {
	jsonapiRuntime := jsonapi.NewRuntime().Instrument("images")

	h.RequestCounter.With(prometheus.Labels{"controller": "gateway", "type": "upload"}).Inc()
	// TODO check header type

	uploadedFileName := uuid.New().String()
	err := minioconnector.UploadFileStream(r.Body, uploadedFileName)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	newImage := &Image{ID: uploadedFileName}
	if err := jsonapiRuntime.MarshalPayload(w, newImage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
