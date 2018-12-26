package api_image

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kmanuel/minioconnector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"io"
	"net/http"
)

type UploadResponse struct {
	FileId string `json:"fileId"`
}

type ImageHandler struct{
	RequestCounter *prometheus.CounterVec
}

func (h *ImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	taskId := mux.Vars(r)["id"]

	log.Info("download request for taskId=", taskId)
	object, err := minioconnector.GetObject(taskId)
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
	h.RequestCounter.With(prometheus.Labels{"controller": "gateway", "type": "upload"}).Inc()
	log.Info("incoming file upload request")
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	file, _, err := r.FormFile("uploadfile")
	defer file.Close()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	uploadedFileName := uuid.New().String()
	err = minioconnector.UploadFileStream(file, uploadedFileName)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var uploadResponse UploadResponse
	uploadResponse.FileId = uploadedFileName

	err = json.NewEncoder(w).Encode(uploadResponse)
	if err != nil {
		log.Error("error writing response")
	}
}