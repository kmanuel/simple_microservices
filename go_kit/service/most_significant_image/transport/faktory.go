package transport

import (
	"bytes"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/service"
	"github.com/opentracing/opentracing-go/log"
)

func CreateFaktoryListenHandler(s service.MostSignificantImageService) worker.Perform {
	return func(ctx worker.Context, args ...interface{}) error {
		task, err := decodeMostSignificantImageTask(args)
		if err != nil {
			_ = ctx.Err()
			return err
		}

		if _, err = s.ExtractMostSignificantImage(task); err != nil {
			log.Error(err)
			_ = ctx.Err()
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
