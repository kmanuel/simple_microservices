package transport

import (
	"bytes"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/service"
)

func CreateFaktoryListenHandler(s service.ScreenshotService) worker.Perform {
	return func(_ worker.Context, args ...interface{}) error {
		task, err := decodeTask(args)
		if err != nil {
			return err
		}
		return s.HandleTask(task)
	}
}

func decodeTask(args []interface{}) (*model.ScreenshotTask, error) {
	task := new(model.ScreenshotTask)
	err := jsonapi.NewRuntime().UnmarshalPayload(bytes.NewBufferString(args[0].(string)), task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
