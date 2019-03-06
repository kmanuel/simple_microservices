package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/simple_microservices/go_kit/gateway/src/model"
	"github.com/kmanuel/simple_microservices/go_kit/gateway/src/service"
	"net/http"
)

func CreateFaktoryHandler(s service.FaktoryService, taskType string) endpoint.Endpoint {
	return func(_ context.Context, task interface{}) (interface{}, error) {
		err := s.PublishToFaktory(taskType, task)
		if err != nil {
			return nil, err
		}
		return task, nil
	}
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return jsonapi.MarshalPayload(w, response)
}

func DecodeCropTask(_ context.Context, r *http.Request) (interface{}, error) {
	task := new(model.CropTask)
	task.ID = uuid.New().String()
	return unmarshal(r, task)
}

func DecodeMostSignificantImageTask(_ context.Context, r *http.Request) (interface{}, error) {
	task := new(model.MostSignificantImageTask)
	task.ID = uuid.New().String()
	return unmarshal(r, task)
}

func DecodeOptimizationTask(_ context.Context, r *http.Request) (interface{}, error) {
	task := new(model.OptimizationTask)
	task.ID = uuid.New().String()
	return unmarshal(r, task)
}

func DecodePortraitTask(_ context.Context, r *http.Request) (interface{}, error) {
	task := new(model.PortraitTask)
	task.ID = uuid.New().String()
	return unmarshal(r, task)
}

func DecodeScreenshotTask(_ context.Context, r *http.Request) (interface{}, error) {
	task := new(model.ScreenshotTask)
	task.ID = uuid.New().String()
	return unmarshal(r, task)
}

func unmarshal(request *http.Request, task interface{}) (interface{}, error) {
	if err := jsonapi.UnmarshalPayload(request.Body, task); err != nil {
		return nil, err
	}
	return task, nil
}

