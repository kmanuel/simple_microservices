package transport

import (
	"bytes"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/service"
)

func CreateFaktoryListenHandler(s service.MostSignificantImageService) worker.Perform {
	return func(_ worker.Context, args ...interface{}) error {
		task, err := decodeMostSignificantImageTask(args)

		_, err = s.ExtractMostSignificantImage(task)

		if err != nil {
			return err
		}

		return nil
	}
}

func decodeMostSignificantImageTask(args []interface{}) (*model.MostSignificantImageTask, error) {
	task := new(model.MostSignificantImageTask)
	err := jsonapi.NewRuntime().UnmarshalPayload(bytes.NewBufferString(args[0].(string)), task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
