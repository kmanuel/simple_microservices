package service

import (
	"github.com/esimov/caire"
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/portrait/model"
	"github.com/prometheus/common/log"
	"os"
)

type OptimizationService interface {
	HandleTask(*model.PortraitTask) error
}

type optimizationServiceImpl struct {
}

func NewOptimizationService() OptimizationService {
	return optimizationServiceImpl{}
}

func (optimizationServiceImpl) HandleTask(task *model.PortraitTask) error {
	downloadedFilePath, err := minioconnector.DownloadFile(task.ImageId)
	if err != nil {
		return err
	}

	outputFilePath, err := extractPortrait(downloadedFilePath, task.Width, task.Height)
	if err != nil {
		return err
	}

	_, err = minioconnector.UploadFileWithName(outputFilePath, task.ID)
	if err != nil {
		return err
	}

	return nil
}

func extractPortrait(inputLocation string, width int, height int) (string, error) {

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

	if err = p.Process(inFile, outFile); err != nil {
		return "", err
	}

	return outputFilePath, nil
}
