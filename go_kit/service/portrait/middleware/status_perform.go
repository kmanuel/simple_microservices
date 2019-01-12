package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/portrait/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/portrait/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/portrait/status_client"
)

type StatusPerformMiddleware struct {
	StatusClient status_client.StatusClient
	Next		 service.OptimizationService
}

func (mw StatusPerformMiddleware) HandleTask(task *model.PortraitTask) error {
	if err := mw.StatusClient.NotifyAboutNew(task.ID); err != nil {
		return err
	}
	if err := mw.Next.HandleTask(task); err != nil {
		_ = mw.StatusClient.NotifyAboutNew(task.ID)
		return err
	}
	return nil
}
