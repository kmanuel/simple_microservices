package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/status_client"
)

type RequestStatusMiddleware struct {
	StatusClient status_client.StatusClient
	Next         service.FaktoryPublishService
}

func (mw RequestStatusMiddleware) PublishTask(task *model.CropTask) error  {
	if err := mw.StatusClient.NotifyAboutNew(task.ID); err != nil {
		return err
	}
	return mw.Next.PublishTask(task)
}
