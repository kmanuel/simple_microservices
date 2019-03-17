package service

import (
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/src/model"
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

type ImageService interface {
	HandleTask(task *model.Task) error
}

func NewCropService(service *minioconnector.MinioService) ImageService {
	return cropServiceImpl{*service}
}

type cropServiceImpl struct {
	minioService minioconnector.MinioService
}

func (c cropServiceImpl) HandleTask(task *model.Task) error {
	imageId := task.ImageId
	width := task.Width
	height := task.Height

	inputImg, err := c.minioService.DownloadFile(imageId)
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

	if _, err = c.minioService.UploadFileWithName(outputFilePath, createFileName(task)); err != nil {
		return err
	}

	log.Info("finished cropping api_image")
	return nil
}

func createFileName(task *model.Task) string {
	inputFileName := strings.Split(task.ImageId, ".")[0]
	taskParams := "height_" + strconv.Itoa(task.Height) + "_width_" + strconv.Itoa(task.Width)
	timestamp := "_" + strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	return inputFileName + "_" + taskParams + timestamp + ".jpg"
}

func (c cropServiceImpl) downloadFile(objectName string) (string, error) {
	return c.minioService.DownloadFile(objectName)
}
