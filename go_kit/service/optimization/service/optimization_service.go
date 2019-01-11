package service

import (
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/model"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/prometheus/common/log"
	"image"
	"image/jpeg"
	"os"
)

type OptimizationService interface {
	HandleTask(*model.OptimizationTask) error
}

type optimizationServiceImpl struct {
}

func NewOptimizationService() OptimizationService {
	return optimizationServiceImpl{}
}

func (optimizationServiceImpl) HandleTask(task *model.OptimizationTask) error {
	downloadedFilePath, err := minioconnector.DownloadFile(task.ImageId)
	if err != nil {
		return err
	}

	outputFilePath, err := optimizeImage(downloadedFilePath, task.Width, task.Height)
	if err != nil {
		return err
	}

	_, err = minioconnector.UploadFileWithName(outputFilePath, task.ID)
	if err != nil {
		return err
	}

	return nil
}

func optimizeImage(inputFile string, width int, height int) (string, error) {
	log.Info("optimizing api_image")
	outputFilePath := "/tmp/" + uuid.New().String() + ".jpg"

	f, _ := os.Open(inputFile)
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

	log.Info("optimized api_image")
	return outputFilePath, nil
}
