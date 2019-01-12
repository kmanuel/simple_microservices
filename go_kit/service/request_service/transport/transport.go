package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/go_kit/service/request_service/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/request_service/service"
	"github.com/prometheus/common/log"
	"net/http"
)

func MakeStatusEndpoint(cs service.RequestStatusService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		log.Info("Status request received")
		list, err := cs.GetTaskStatusList()
		if err != nil {
			log.Error(err)
			return nil, err
		}
		return &list, nil
	}
}

func MakeStatusChangeEndpoint(cs service.ChangeStatusService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		log.Info("Status change request receives")
		task := request.(*model.TaskStatus)
		return cs.SaveOrUpdate(task)
	}
}

func DecodeTaskStatus(_ context.Context, r *http.Request) (interface{}, error) {
	task := new(model.TaskStatus)
	if err := jsonapi.UnmarshalPayload(r.Body, task); err != nil {
		return nil, err
	}
	return task, nil
}

func DecodeListRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return jsonapi.MarshalPayload(w, response)
}
