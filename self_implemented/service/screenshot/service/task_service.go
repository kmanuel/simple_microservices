package service

import (
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/service/screenshot/model"
	"os/exec"
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
	outputFilePath, err := takeScreenShot(t.Url)
	if err != nil {
		return err
	}

	_, err = minioconnector.UploadFileWithName(outputFilePath, t.ID)
	if err != nil {
		return err
	}

	return nil
}

func takeScreenShot(url string) (string, error) {
	chromeUserAgent := "Mozilla/5.0 (Windows NT 6.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	phantomJSBin := "./lib/bin/phantomjs"
	jsPath := "./lib/js/screenshot.js"
	logFile := "screenshot_service.log"

	outputFilePath := "/tmp/output" + uuid.New().String() + ".jpg"

	cmd := exec.Command(phantomJSBin, jsPath, url, outputFilePath, logFile, chromeUserAgent)

	if err := cmd.Run(); nil != err {
		return "", err
	}

	return outputFilePath, nil
}