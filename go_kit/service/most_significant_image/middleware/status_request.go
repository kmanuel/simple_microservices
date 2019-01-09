package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/status_client"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/service"
)

type StatusRequestMiddleware struct {
	StatusClient status_client.StatusClient
	Next		 service.FaktoryPublishService
}

func (mw StatusRequestMiddleware) PublishTask(task *model.MostSignificantImageTask) error {
	if err := mw.StatusClient.NotifyAboutNew(task.ID); err != nil {
		return err
	}
	return mw.Next.PublishTask(task)
}
