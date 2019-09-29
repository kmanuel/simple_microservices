package function

import (
	"bytes"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/prometheus/common/log"
	"image"
	"image/jpeg"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `jsonapi:"primary,crop_task"`
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
		"crop")
}

func handleTask(task *Task, minioService minioconnector.MinioService) error {
	imageId := task.ImageId
	width := task.Width
	height := task.Height

	inputImg, err := minioService.DownloadFile(imageId)
	if err != nil {
		return err
	}

	log.Info("starting to crop api_image")

	outputFilePath := "/tmp/downloaded" + uuid.New().String() + ".jpg"

	f, _ := os.Open(inputImg)
	img, _, _ := image.Decode(f)
	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	topCrop, _ := analyzer.FindBestCrop(img, width, height)
	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	croppedImg := img.(SubImager).SubImage(topCrop)
	f, err = os.Create(outputFilePath)
	defer f.Close()
	if err != nil {
		return err
	}
	if err = jpeg.Encode(f, croppedImg, nil); err != nil {
		return err
	}

	if _, err = minioService.UploadFileWithName(outputFilePath, createFileName(task)); err != nil {
		return err
	}

	log.Info("finished cropping api_image")
	return nil
}

func createFileName(task *Task) string {
	inputFileName := strings.Split(task.ImageId, ".")[0]
	taskParams := "height_" + strconv.Itoa(task.Height) + "_width_" + strconv.Itoa(task.Width)
	timestamp := "_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	return inputFileName + "_" + taskParams + timestamp + ".jpg"
}
