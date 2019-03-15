package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kmanuel/simple_microservices/go_kit/service/gateway/src/service"
	"net/http"
)

func CreateRestHandler(s service.FaktoryService) endpoint.Endpoint {
	return func(_ context.Context, task interface{}) (interface{}, error) {
		info, err := s.Info()
		if err != nil {
			return nil, err
		}
		return info, nil
	}
}

func DecodeInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}
