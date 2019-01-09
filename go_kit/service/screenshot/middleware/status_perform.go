package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/status_client"
)

type StatusPerformMiddleware struct {
	StatusClient status_client.StatusClient
	Next		 service.ScreenshotService
}

func (mw StatusPerformMiddleware) HandleTask(task *model.ScreenshotTask) error {
	if err := mw.StatusClient.NotifyAboutNew(task.ID); err != nil {
		return err
	}
	return mw.Next.HandleTask(task)
}
