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

func handleTask(task *Task) error {
	chromeUserAgent := "Mozilla/5.0 (Windows NT 6.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	phantomJSBin := "/home/app/lib/bin/phantomjs"
	jsPath := "/home/app/lib/js/screenshot.js"
	logFile := "screenshot_service.log"

	outputFilePath := "/tmp/output" + uuid.New().String() + ".jpg"

	cmd := exec.Command(phantomJSBin, jsPath, task.Url, outputFilePath, logFile, chromeUserAgent)

	if err := cmd.Run(); nil != err {
		return err
	}

	if _, err := minioconnector.UploadFileWithName(outputFilePath, createFileName(task)); err != nil {
		return err
	}

	return nil
}

func createFileName(task *Task) string {
	inputFileName := strings.Replace(task.Url, ".", "_", -1)
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	return inputFileName + "_" + timestamp + "_screenshot.jpg"
}
