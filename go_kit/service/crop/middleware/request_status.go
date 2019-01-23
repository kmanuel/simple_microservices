package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/service"
)

type RequestStatusMiddleware struct {
	statusClient service.StatusClient
	next         service.ImageService
}

func NewRequestStatusMiddleware(statusClient service.StatusClient, next service.ImageService) service.ImageService {
	return RequestStatusMiddleware{statusClient: statusClient, next: next}
}

func (mw RequestStatusMiddleware) HandleTask(task *model.Task) error {
	if err := mw.statusClient.NotifyAboutProcessing(task.ID); err != nil {
		return err
	}

	err := mw.next.HandleTask(task)

	if err != nil {
		if err := mw.statusClient.NotifyAboutFailure(task.ID); err != nil {
			return err
		}
		return err
	}
	if err := mw.statusClient.NotifyAboutCompletion(task.ID); err != nil {
		return err
	}
	return nil
}
