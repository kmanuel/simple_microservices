package transport

import (
	"bytes"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/service"
	"github.com/opentracing/opentracing-go/log"
)

func CreateFaktoryListenHandler(s service.ImageService) worker.Perform {
	return func(ctx worker.Context, args ...interface{}) error {
		task, err := decodeTask(args)
		if err != nil {
			_ = ctx.Err()
			return err
		}
		if err = s.HandleTask(task); err != nil {
			log.Error(err)
			_ = ctx.Err()
			return err
		}
		ctx.Done()
		return nil
	}
}

func decodeTask(args []interface{}) (*model.Task, error) {
	task := new(model.Task)
	err := jsonapi.NewRuntime().UnmarshalPayload(bytes.NewBufferString(args[0].(string)), task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
