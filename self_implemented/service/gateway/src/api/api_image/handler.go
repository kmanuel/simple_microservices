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
	DispatchCounter *prometheus.CounterVec
	MinioService minioconnector.MinioService
}

func (h *ImageHandler) DownloadImage(w http.ResponseWriter, r *http.Request) {
	log.Info("download api_image request")
	imageId := mux.Vars(r)["id"]

	log.Info("download request for imageId=", imageId)
	object, err := h.MinioService.GetObject(imageId)
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

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	h.DispatchCounter.With(prometheus.Labels{"type": "upload"}).Inc()

	uploadedFileName := uuid.New().String()
	err := h.MinioService.UploadFileStream(r.Body, uploadedFileName)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	newImage := &Image{ID: uploadedFileName}
	if err := jsonapi.MarshalPayload(w, newImage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
