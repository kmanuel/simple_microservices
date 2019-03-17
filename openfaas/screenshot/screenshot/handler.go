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
	ID  string `jsonapi:"primary,screenshot_task"`
	Url string `jsonapi:"attr,url"`
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
		"screenshot")
}

func handleTask(task *Task, minioService minioconnector.MinioService) error {
	chromeUserAgent := "Mozilla/5.0 (Windows NT 6.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	phantomJSBin := "/home/app/lib/bin/phantomjs"
	jsPath := "/home/app/lib/js/screenshot.js"
	logFile := "screenshot_service.log"

	outputFilePath := "/tmp/output" + uuid.New().String() + ".jpg"

	cmd := exec.Command(phantomJSBin, jsPath, task.Url, outputFilePath, logFile, chromeUserAgent)

	if err := cmd.Run(); nil != err {
		return err
	}

	if _, err := minioService.UploadFileWithName(outputFilePath, createFileName(task)); err != nil {
		return err
	}

	return nil
}

func createFileName(task *Task) string {
	inputFileName := strings.Replace(task.Url, "http://", "", -1)
	inputFileName = strings.Replace(inputFileName, "https://", "", -1)
	inputFileName = strings.Replace(inputFileName, ".", "_", -1)
	inputFileName = strings.Replace(inputFileName, "/", "_", -1)
	timestamp := "_" + strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	return inputFileName + timestamp + ".jpg"
}
