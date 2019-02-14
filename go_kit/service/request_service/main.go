package main

import (
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/request_service/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/request_service/transport"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	initMinio()

	var ss service.RequestStatusService
	ss = service.RequestStatusServiceImpl{}

	requestListHandler := httptransport.NewServer(
		transport.MakeStatusEndpoint(ss),
		transport.DecodeListRequest,
		transport.EncodeResponse,
	)
	http.Handle("/tasks", requestListHandler)

	var changeStatusService service.ChangeStatusService
	changeStatusService = service.RequestStatusServiceImpl{}
	changeStatusHandler := httptransport.NewServer(
		transport.MakeStatusChangeEndpoint(changeStatusService),
		transport.DecodeTaskStatus,
		transport.EncodeResponse,
	)
	http.Handle("/", changeStatusHandler)

	fmt.Println(http.ListenAndServe(":8080", nil))
}

func initMinio() {
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))
}

