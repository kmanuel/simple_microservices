package api_image

import (
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kmanuel/minioconnector"
	"github.com/prometheus/common/log"
	"io"
	"net/http"
)

type ImageHandler struct {}

func (h ImageHandler) ServeUploadHTTP(w http.ResponseWriter, r *http.Request) {
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

	uploadedFileName := uuid.New().String()
	err := minioconnector.UploadFileStream(r.Body, uploadedFileName)
	if err != nil {
		log.Error(err)
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

func (h ImageHandler) ServeDownload(w http.ResponseWriter, r *http.Request) {
	imageId := mux.Vars(r)["id"]

	object, err := minioconnector.GetObject(imageId)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}
	_, err = io.Copy(w, object)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}
}

func (h *ImageHandler) downloadImage(w http.ResponseWriter, r *http.Request) {
	imageId := mux.Vars(r)["id"]

	object, err := minioconnector.GetObject(imageId)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}
	_, err = io.Copy(w, object)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}
}
