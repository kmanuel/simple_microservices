package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/middlewares/prometheus"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/transport"
	"net/http"
	"os"
)

type loggingMiddleware struct {
	logger log.Logger
	next   service.CropService
}

func (mw loggingMiddleware) CropImage(inputImg string, width int, height int) (output string, err error) {
	_ = mw.logger.Log("method", "cropImage")
	output, err = mw.next.CropImage(inputImg, width, height)
	return
}

func initMinio() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))
}

func main() {
	initMinio()

	go prometheus.Init()

	logger := log.NewLogfmtLogger(os.Stdout)

	var cs service.CropService
	cs = service.CropServiceImpl{}
	cs = loggingMiddleware{logger, cs}
	cs = prometheus.Middleware{Next: cs}

	cropHandler := httptransport.NewServer(
		transport.MakeCropEndpoint(cs),
		transport.DecodeCropRequest,
		transport.EncodeResponse,
	)

	http.Handle("/", cropHandler)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

