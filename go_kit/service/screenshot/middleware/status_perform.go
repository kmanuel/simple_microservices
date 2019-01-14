package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/status_client"
	"github.com/prometheus/common/log"
)

type StatusPerformMiddleware struct {
	StatusClient status_client.StatusClient
	Next		 service.ScreenshotService
}

func (mw StatusPerformMiddleware) HandleTask(task *model.ScreenshotTask) error {
	if err := mw.StatusClient.NotifyAboutNew(task.ID); err != nil {
		return err
	}
	if err := mw.Next.HandleTask(task); err != nil {
		log.Error(err)
		_ = mw.StatusClient.NotifyAboutFailure(task.ID)
		return err
	}
	if err := mw.StatusClient.NotifyAboutCompletion(task.ID); err != nil {
		log.Error("failed to notify request_status service about failure")
		return err
	}
	return nil
}
