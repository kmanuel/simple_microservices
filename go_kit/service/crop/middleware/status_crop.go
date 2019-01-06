package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/status_client"
)

type StatusCropMiddleware struct {
	StatusClient status_client.StatusClient
	Next         service.CropService
}

func (mw StatusCropMiddleware) CropImage(task *model.CropTask) (string, error)  {
	if err := mw.StatusClient.NotifyAboutProcessing(task.ID); err != nil {
		return "", err
	}
	res, err := mw.Next.CropImage(task)
	if err != nil {
		if err := mw.StatusClient.NotifyAboutFailure(task.ID); err != nil {
			return "", err
		}
	} else {
		if err := mw.StatusClient.NotifyAboutCompletion(task.ID); err != nil {
			return "", err
		}
	}
	return res, err
}
