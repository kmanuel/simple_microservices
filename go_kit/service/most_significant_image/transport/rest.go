package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/service"
	"net/http"
)

func MakeMostSignificantImageEndpoint(s service.MostSignificantImageService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		task := request.(*model.MostSignificantImageTask)
		task.ID = uuid.New().String()
		_, err := s.ExtractMostSignificantImage(task)
		if err != nil {
			return nil, err
		}
		return task, nil
	}
}

func DecodeMostSignificantImageTask(_ context.Context, r *http.Request) (interface{}, error) {
	task := new(model.MostSignificantImageTask)
	if err := jsonapi.UnmarshalPayload(r.Body, task); err != nil {
		return nil, err
	}
	return task, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return jsonapi.MarshalPayload(w, response)
}
