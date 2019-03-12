package function

import (
	"bytes"
	"fmt"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/prometheus/common/log"
	"os"
	"os/exec"
)

type Task struct {
	ID      string `jsonapi:"primary,optimization_task"`
	ImageId string `jsonapi:"attr,image_id"`
	Width   int    `jsonapi:"attr,width"`
	Height  int    `jsonapi:"attr,height"`
}

// Handle a serverless request
func Handle(req []byte) string {

	initMinio()

	task := new(Task)
	task.ID = uuid.New().String()

	jsonapi.UnmarshalPayload(bytes.NewReader(req), task)

	err := handleTask(task)
	if err != nil {
		log.Error("wowowow, error", err)
	}

	return fmt.Sprintf("oki")
}

func initMinio() {
	minioHost := os.Getenv("MINIO_HOST")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecret := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("BUCKET_NAME")

	log.Errorf("initializing minio with host=%s accessKey=%s secret=%s bucketName=%s", minioHost, minioAccessKey, minioSecret, bucketName)
	minioconnector.Init(
		minioHost,
		minioAccessKey,
		minioSecret,
		bucketName)
}

func handleTask(t *Task) error {

	downloadedFilePath, err := minioconnector.DownloadFile(t.ImageId)
	if err != nil {
		return err
	}

	outputFilePath, err := optimizeImage(downloadedFilePath)
	if err != nil {
		return err
	}

	_, err = minioconnector.UploadFileWithName(outputFilePath, t.ID)
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

