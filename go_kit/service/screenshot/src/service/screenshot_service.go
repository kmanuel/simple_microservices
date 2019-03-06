package service

import (
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/src/model"
	"os/exec"
)

type ImageService interface {
	HandleTask(*model.Task) error
}

type screenshotServiceImpl struct {
	minioService minioconnector.MinioService
}

func NewScreenshotService(minioService *minioconnector.MinioService) ImageService {
	return screenshotServiceImpl{*minioService}
}

func (s screenshotServiceImpl) HandleTask(task *model.Task) error {
	chromeUserAgent := "Mozilla/5.0 (Windows NT 6.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	phantomJSBin := "./lib/bin/phantomjs"
	jsPath := "./lib/js/screenshot.js"
	logFile := "screenshot_service.log"

	outputFilePath := "/tmp/output" + uuid.New().String() + ".jpg"

	cmd := exec.Command(phantomJSBin, jsPath, task.Url, outputFilePath, logFile, chromeUserAgent)

	if err := cmd.Run(); nil != err {
		return err
	}

	if _, err := s.minioService.UploadFileWithName(outputFilePath, task.ID); err != nil {
		return err
	}

	return nil
}
