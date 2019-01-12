package service

import (
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/portrait/model"
	"github.com/prometheus/common/log"
	"os/exec"
	"strconv"
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

	cmd := exec.Command(
		"caire",
		"-in", inputLocation,
		"-out", outputFilePath,
		"-width="+strconv.Itoa(width),
		"-height="+strconv.Itoa(height),
		"-perc=1",
		"-square=0",
		"-scale=1",
		"-blur=0",
		"-sobel=0",
		"-debug=0",
		"-face=1",
		"-cc=./data/facefinder",
	)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	log.Info("extracted portrait")
	return outputFilePath, nil
}
