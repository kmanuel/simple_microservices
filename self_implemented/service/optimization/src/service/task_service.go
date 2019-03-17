package service

import (
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/service/optimization/src/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

	_, err = h.minioService.UploadFileWithName(outputFilePath, createFileName(t))
	if err != nil {
		return err
	}

	return nil
}

func createFileName(task *model.Task) string {
	urlFileName := strings.Split(task.ImageId, ".")[0]
	timestamp := "_" + strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	return urlFileName + timestamp + ".jpg"
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
