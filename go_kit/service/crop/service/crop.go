package service

import (
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/model"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/prometheus/common/log"
	"image"
	"image/jpeg"
	"os"
)

type ImageService interface {
	HandleTask(task *model.Task) error
}

func NewCropService() ImageService {
	return cropServiceImpl{}
}

type cropServiceImpl struct{}

func (cropServiceImpl) HandleTask(task *model.Task) error {
	imageId := task.ImageId
	width := task.Width
	height := task.Height

	inputImg, err := downloadFile(imageId)
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

	if _, err = minioconnector.UploadFileWithName(outputFilePath, task.ID); err != nil {
		return err
	}

	log.Info("finished cropping api_image")
	return nil
}

func downloadFile(objectName string) (string, error) {
	return minioconnector.DownloadFile(objectName)
}
