package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/go_kit/service/request_service/service"
	"net/http"
)

func MakeStatusEndpoint(cs service.RequestStatusService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		list, err := cs.GetTaskStatusList()
		if err != nil {
			panic(err) // TODO
		}
		return &list, nil
	}
}

func DecodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return jsonapi.MarshalPayload(w, response)
}
