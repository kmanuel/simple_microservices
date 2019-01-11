package transport

import (
	"bytes"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/service"
)

func CreateFaktoryListenHandler(s service.OptimizationService) worker.Perform {
	return func(_ worker.Context, args ...interface{}) error {
		task, err := decodeTask(args)
		if err != nil {
			return err
		}
		return s.HandleTask(task)
	}
}

func decodeTask(args []interface{}) (*model.OptimizationTask, error) {
	task := new(model.OptimizationTask)
	err := jsonapi.NewRuntime().UnmarshalPayload(bytes.NewBufferString(args[0].(string)), task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
