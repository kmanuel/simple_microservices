package function

import (
	"bytes"
	"fmt"
	"github.com/esimov/caire"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/prometheus/common/log"
	"os"
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

	outputFilePath, err := ExtractPortrait(downloadedFilePath, t.Width, t.Height)
	if err != nil {
		return err
	}

	_, err = minioconnector.UploadFileWithName(outputFilePath, createFileName(t))
	if err != nil {
		return err
	}

	return nil
}

func createFileName(task *Task) string {
	inputFileName := strings.Split(task.ImageId, ".")[0]
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	taskParams := "height_" + strconv.Itoa(task.Height) + "_width_" + strconv.Itoa(task.Width)
	return inputFileName + "_" + timestamp + "_portrait_" + taskParams + ".jpg"
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
