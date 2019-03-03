package service

import (
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/service/optimization/model"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"image"
	"image/jpeg"
	"os"
)

type TaskService interface {
	Handle(t *model.Task) error
}

type taskService struct {
	counter      *prometheus.CounterVec
	taskType     string
	minioService minioconnector.MinioService
}

func NewTaskService(counter *prometheus.CounterVec, taskType string, minioService *minioconnector.MinioService) TaskService {
	return taskService{counter, taskType, *minioService}
}

func (h taskService) Handle(t *model.Task) error {
	h.counter.With(prometheus.Labels{"type": h.taskType}).Inc()

	downloadedFilePath, err := h.minioService.DownloadFile(t.ImageId)
	if err != nil {
		return err
	}

	outputFilePath, err := optimizeImage(downloadedFilePath, t.Width, t.Height)
	if err != nil {
		return err
	}

	_, err = h.minioService.UploadFileWithName(outputFilePath, t.ID)
	if err != nil {
		return err
	}

	return nil
}

func optimizeImage(inputFile string, width int, height int) (string, error) {
	log.Info("optimizing api_image")
	outputFilePath := "/tmp/" + uuid.New().String() + ".jpg"

	f, _ := os.Open(inputFile)
	img, _, _ := image.Decode(f)
	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	topCrop, _ := analyzer.FindBestCrop(img, width, height)
	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	croppedImg := img.(SubImager).SubImage(topCrop)
	f, err := os.Create(outputFilePath)
	defer f.Close()
	if err != nil {
		return "", err
	}
	err = jpeg.Encode(f, croppedImg, nil)
	if err != nil {
		return "", err
	}

	log.Info("optimized api_image")
	return outputFilePath, nil
}
