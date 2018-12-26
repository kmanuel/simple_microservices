package image

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/kmanuel/minioconnector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"io"
	"net/http"
)

type UploadResponse struct {
	FileId string `json:"fileId"`
}

type ImageController struct {}

var requestsCounter *prometheus.CounterVec

func NewImageController(requestsCounterArg *prometheus.CounterVec) *ImageController {
	requestsCounter = requestsCounterArg
	return &ImageController{}
}

func (s *ImageController) HandleUpload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requestsCounter.With(prometheus.Labels{"controller":"gateway", "type": "upload"}).Inc()
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

func (s *ImageController) HandleDownload(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	requestsCounter.With(prometheus.Labels{"controller":"gateway", "type": "download"}).Inc()
	taskId := ps.ByName("taskId")
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