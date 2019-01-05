package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/service"
	"net/http"
)

func MakeCropEndpoint(cs service.CropService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(*model.CropTask)
		req.ID = uuid.New().String()
		_, err := cs.CropImage(req.ImageId, req.Width, req.Height)
		if err != nil {
			return nil, err
		}
		return req, nil
	}
}

func DecodeCropRequest(_ context.Context, r *http.Request) (interface{}, error) {
	task := new(model.CropTask)
	if err := jsonapi.UnmarshalPayload(r.Body, task); err != nil {
		return nil, err
	}
	return task, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return jsonapi.MarshalPayload(w, response)
}

