package service

import (
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/service/crop/model"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/prometheus/common/log"
	"image"
	"image/jpeg"
	"os"
)

type TaskService interface {
	Handle(t *model.Task) error
}

type taskService struct {
}

func NewTaskService() TaskService {
	return taskService{}
}

func (h taskService) Handle(t *model.Task) error {
	downloadedFilePath, err := downloadFile(t.ImageId)
	if err != nil {
		return err
	}
	croppedFilePath, err := cropImage(downloadedFilePath, t.Width, t.Height)
	if err != nil {
		return err
	}
	_, err = minioconnector.UploadFileWithName(croppedFilePath, t.ID)
	return err
}

func downloadFile(objectName string) (string, error) {
	return minioconnector.DownloadFile(objectName)
}

func cropImage(inputImg string, width int, height int) (string, error) {
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
	f, err := os.Create(outputFilePath)
	defer f.Close()
	if err != nil {
		return "", err
	}
	err = jpeg.Encode(f, croppedImg, nil)
	if err != nil {
		return "", err
	}

	log.Info("finished cropping api_image")
	return outputFilePath, nil
}
