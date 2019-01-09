package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/status_client"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/service"
)

type StatusPerformMiddleware struct {
	StatusClient status_client.StatusClient
	Next		 service.MostSignificantImageService
}

func (mw StatusPerformMiddleware) ExtractMostSignificantImage(task *model.MostSignificantImageTask) (outputImagePath string, err error) {
	if err := mw.StatusClient.NotifyAboutNew(task.ID); err != nil {
		return "nil", err
	}
	return mw.Next.ExtractMostSignificantImage(task)
}
