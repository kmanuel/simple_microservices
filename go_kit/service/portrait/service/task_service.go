package service

import (
	"github.com/kmanuel/caire"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/portrait/model"
	"github.com/prometheus/common/log"
	"os"
)

type ImageService interface {
	HandleTask(*model.Task) error
}

type optimizationServiceImpl struct {
}

func NewOptimizationService() ImageService {
	return optimizationServiceImpl{}
}

func (optimizationServiceImpl) HandleTask(task *model.Task) error {
	downloadedFilePath, err := minioconnector.DownloadFile(task.ImageId)
	if err != nil {
		return err
	}

	outputFilePath, err := extractPortrait(task.ID, downloadedFilePath, task.Width, task.Height)
	if err != nil {
		return err
	}

	_, err = minioconnector.UploadFileWithName(outputFilePath, task.ID)
	if err != nil {
		return err
	}

	return nil
}

func extractPortrait(taskId string, inputLocation string, width int, height int) (string, error) {

	log.Info("extracting portrait")

	outputFilePath := "/tmp/" + taskId + ".jpg"

	p := caire.Processor{
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
		return "", err
	}

	outFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY, 0755)
	defer outFile.Close()
	if err != nil {
		log.Fatalf("Unable to open output file: %v", err)
		return "", err
	}

	log.Info("processing file")
	if err = p.Process(inFile, outFile); err != nil {
		log.Error("foo")
		return "", err
	}

	log.Info("extracted portrait")

	return outputFilePath, nil
}
