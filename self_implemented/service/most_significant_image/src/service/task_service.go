package service

import (
	"github.com/advancedlogic/GoOse"
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/service/most_significant_image/src/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const outputImageLocation = "/tmp/"

type TaskService interface {
	Handle(t *model.Task) error
}

type taskService struct {
	counter      *prometheus.CounterVec
	taskType     string
	minioService minioconnector.MinioService
}

func NewTaskService(counter *prometheus.CounterVec,
	taskType string,
	minioService *minioconnector.MinioService) TaskService {
	return taskService{counter, taskType, *minioService}
}

func (h taskService) Handle(t *model.Task) error {
	h.counter.With(prometheus.Labels{"type": h.taskType}).Inc()

	task := t

	outputFile := outputImageLocation + uuid.New().String() + ".jpg"

	err := ExtractMostSignificantImage(task.Url, outputFile)
	if err != nil {
		log.Error("extracting of image failed", err)
		return err
	}

	fileName := h.createFileName(t)
	_, err = h.minioService.UploadFileWithName(outputFile, fileName)
	if err != nil {
		log.Error("upload failed", err)
		return err
	}

	return nil
}

func (h taskService) createFileName(task *model.Task) string {
	inputFileName := strings.Replace(task.Url, ".", "_", -1)
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	return inputFileName + "_" + timestamp + "_" + h.taskType + ".jpg"
}

func ExtractMostSignificantImage(inputUrl string, outputFile string) error {
	g := goose.New()
	article, err := g.ExtractFromURL(inputUrl)
	if err != nil {
		return err
	}
	topImageUrl := article.TopImage
	err = DownloadImage(topImageUrl, outputFile)
	return err
}

func DownloadImage(url string, outputFile string) error {
	filepath := outputFile

	out, err := os.Create(filepath)
	defer out.Close()
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
