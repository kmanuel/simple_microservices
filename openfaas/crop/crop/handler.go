package function

import (
	"bytes"
	"fmt"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/prometheus/common/log"
	"image"
	"image/jpeg"
	"os"
)

type Task struct {
	ID      string `jsonapi:"primary,crop_task"`
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

func handleTask(task *Task) error {
	imageId := task.ImageId
	width := task.Width
	height := task.Height

	inputImg, err := downloadFile(imageId)
	if err != nil {
		log.Error("meh, downloadFile failed")
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

	if _, err = minioconnector.UploadFileWithName(outputFilePath, task.ID); err != nil {
		return err
	}

	log.Info("finished cropping api_image")
	return nil
}

func downloadFile(objectName string) (string, error) {
	return minioconnector.DownloadFile(objectName)
}
