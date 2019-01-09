package service

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/status_client"
)

type ScreenshotService interface {
	HandleTask(*model.ScreenshotTask) error
}

type screenshotTaskImpl struct {
	statusClient status_client.StatusClient
}

func NewScreenshotService(s status_client.StatusClient) ScreenshotService {
	return screenshotTaskImpl{s}
}

func (screenshotTaskImpl) HandleTask(*model.ScreenshotTask) error {
	panic("implement me")
}
