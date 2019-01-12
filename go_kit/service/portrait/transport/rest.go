package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/simple_microservices/go_kit/service/portrait/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/portrait/service"
	"net/http"
)

func CreateRestHandler(s service.FaktoryPublishService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		task := request.(*model.PortraitTask)
		task.ID = uuid.New().String()
		err := s.PublishTask(task)
		if err != nil {
			return nil, err
		}
		return task, nil
	}
}

func DecodeScreenshotTask(_ context.Context, r *http.Request) (interface{}, error) {
	task := new(model.PortraitTask)
	if err := jsonapi.UnmarshalPayload(r.Body, task); err != nil {
		return nil, err
	}
	return task, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return jsonapi.MarshalPayload(w, response)
}
