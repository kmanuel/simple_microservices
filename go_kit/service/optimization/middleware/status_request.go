package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/status_client"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/service"
)

type StatusRequestMiddleware struct {
	StatusClient status_client.StatusClient
	Next		 service.FaktoryPublishService
}

func (mw StatusRequestMiddleware) PublishTask(task *model.OptimizationTask) error {
	if err := mw.StatusClient.NotifyAboutNew(task.ID); err != nil {
		return err
	}
	return mw.Next.PublishTask(task)
}
