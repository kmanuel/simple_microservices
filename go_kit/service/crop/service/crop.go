package service

import (
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/prometheus/common/log"
	"image"
	"image/jpeg"
	"os"
)

type CropService interface {
	CropImage(inputImg string, width int, height int) (string, error)
}

type CropServiceImpl struct{}

func downloadFile(objectName string) (string, error) {
	return minioconnector.DownloadFile(objectName)
}

func (CropServiceImpl) CropImage(imageId string, width int, height int) (string, error) {
	inputImg, err := downloadFile(imageId)
	if err != nil {
		return "", err
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
		return "", err
	}
	err = jpeg.Encode(f, croppedImg, nil)
	if err != nil {
		return "", err
	}

	_, err = minioconnector.UploadFileWithName(outputFilePath, imageId)

	log.Info("finished cropping api_image")
	return outputFilePath, nil
}
