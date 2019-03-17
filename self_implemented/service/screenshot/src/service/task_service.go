package service

import (
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/service/screenshot/src/model"
	"github.com/prometheus/client_golang/prometheus"
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

	outputFilePath, err := takeScreenShot(t.Url)
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
	inputFileName := strings.Replace(task.Url, "http://", "", -1)
	inputFileName = strings.Replace(inputFileName, "https://", "", -1)
	inputFileName = strings.Replace(inputFileName, ".", "_", -1)
	inputFileName = strings.Replace(inputFileName, "/", "_", -1)
	timestamp := "_" + strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	return inputFileName + timestamp + ".jpg"
}

func takeScreenShot(url string) (string, error) {
	chromeUserAgent := "Mozilla/5.0 (Windows NT 6.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	phantomJSBin := "./lib/bin/phantomjs"
	jsPath := "./lib/js/screenshot.js"
	logFile := "screenshot_service.log"

	outputFilePath := "/tmp/output" + uuid.New().String() + ".jpg"

	cmd := exec.Command(phantomJSBin, jsPath, url, outputFilePath, logFile, chromeUserAgent)

	if err := cmd.Run(); nil != err {
		return "", err
	}

	return outputFilePath, nil
}
