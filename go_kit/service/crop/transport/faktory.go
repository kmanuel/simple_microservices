package transport

import (
	"bytes"
	"fmt"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/service"
)

func CreateFaktoryHandler(cs service.CropService) worker.Perform {
	return func(_ worker.Context, args ...interface{}) error {
		task, err := decodeCropTask(args)
		if err != nil {
			return err
		}

		if _, err = cs.CropImage(task); err != nil {
			fmt.Println("error while cropping image", err)
			return err
		}

		return err
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
