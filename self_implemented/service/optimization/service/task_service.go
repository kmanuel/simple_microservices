package service

import (
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/service/optimization/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"os/exec"
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

	outputFilePath, err := optimizeImage(downloadedFilePath)
	if err != nil {
		return err
	}

	_, err = h.minioService.UploadFileWithName(outputFilePath, t.ID)
	if err != nil {
		return err
	}

	return nil
}

func optimizeImage(inputFile string) (string, error) {
	log.Info("optimizing api_image")
	cmd := exec.Command("image_optim", inputFile)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	log.Info("optimized api_image")
	return inputFile, nil
}
