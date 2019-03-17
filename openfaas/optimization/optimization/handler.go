package function

import (
	"bytes"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/prometheus/common/log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `jsonapi:"primary,optimization_task"`
	ImageId string `jsonapi:"attr,image_id"`
	Width   int    `jsonapi:"attr,width"`
	Height  int    `jsonapi:"attr,height"`
}

// Handle a serverless request
func Handle(req []byte) string {

	minioService := initMinio()

	task := new(Task)
	task.ID = uuid.New().String()

	err := jsonapi.UnmarshalPayload(bytes.NewReader(req), task)
	if err != nil {
		panic(err)
	}

	err = handleTask(task, *minioService)
	if err != nil {
		panic(err)
	}

	return ""
}

func initMinio() *minioconnector.MinioService {
	minioHost := os.Getenv("MINIO_HOST")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecret := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("INPUT_BUCKET_NAME")

	log.Errorf("initializing minio with host=%s accessKey=%s secret=%s bucketName=%s", minioHost, minioAccessKey, minioSecret, bucketName)
	return minioconnector.NewMinioService(
		minioHost,
		minioAccessKey,
		minioSecret,
		bucketName,
		"optimization")
}

func handleTask(t *Task, minioService minioconnector.MinioService) error {

	downloadedFilePath, err := minioService.DownloadFile(t.ImageId)
	if err != nil {
		return err
	}

	outputFilePath, err := optimizeImage(downloadedFilePath)
	if err != nil {
		return err
	}

	_, err = minioService.UploadFileWithName(outputFilePath, createFileName(t))
	if err != nil {
		return err
	}

	return nil
}

func createFileName(task *Task) string {
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
