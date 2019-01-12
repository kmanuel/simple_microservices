package transport

import (
	"bytes"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/service"
)

func CreateFaktoryHandler(cs service.CropService) worker.Perform {
	return func(ctx worker.Context, args ...interface{}) error {
		task, err := decodeCropTask(args)
		if err != nil {
			_ = ctx.Err()
			return err
		}

		if _, err = cs.CropImage(task); err != nil {
			_ = ctx.Err()
			return err
		}

		return nil
	}
}

func decodeCropTask(args []interface{}) (*model.CropTask, error) {
	task := new(model.CropTask)
	err := jsonapi.NewRuntime().UnmarshalPayload(bytes.NewBufferString(args[0].(string)), task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
