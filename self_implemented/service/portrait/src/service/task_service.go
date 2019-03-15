package service

import (
	"github.com/google/uuid"
	"github.com/kmanuel/caire"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/service/portrait/src/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"os"
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

	outputFilePath, err := ExtractPortrait(downloadedFilePath, t.Width, t.Height)
	if err != nil {
		return err
	}

	outputFileName := h.createFileName(t)
	_, err = h.minioService.UploadFileWithName(outputFilePath, outputFileName)
	if err != nil {
		return err
	}

	return nil
}

func (h taskService) createFileName(task *model.Task) string {
	inputFileName := strings.Split(task.ImageId, ".")[0]
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	taskParams := "height_" + strconv.Itoa(task.Height) + "_width_" + strconv.Itoa(task.Width)
	return inputFileName + "_" + timestamp + "_" + h.taskType + "_" + taskParams + ".jpg"
}

func ExtractPortrait(inputLocation string, width int, height int) (string, error) {

	log.Info("extracting portrait")

	outputFilePath := "/tmp/" + uuid.New().String() + ".jpg"

	p := &caire.Processor{
		BlurRadius:     0,
		SobelThreshold: 0,
		NewWidth:       width,
		NewHeight:      height,
		Percentage:     true,
		Square:         false,
		Debug:          false,
		Scale:          true,
		FaceDetect:     true,
		Classifier:     "./data/facefinder",
	}

	inFile, err := os.Open(inputLocation)
	defer inFile.Close()
	if err != nil {
		log.Fatalf("Unable to open source file: %v", err)
	}

	outFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY, 0755)
	defer outFile.Close()
	if err != nil {
		log.Fatalf("Unable to open output file: %v", err)
	}

	err = p.Process(inFile, outFile)
	if err != nil {
		return "", err
	}

	log.Info("extracted portrait")
	return outputFilePath, nil
}
