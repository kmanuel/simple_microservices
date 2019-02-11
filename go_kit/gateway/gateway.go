package main

import (
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/kmanuel/simple_microservices/go_kit/gateway/transport"

	//"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	//"github.com/kmanuel/simple_microservices/go_kit/gateway/api_image"
	"github.com/kmanuel/simple_microservices/go_kit/gateway/service"
	"github.com/prometheus/common/log"
	"net/http"
	"os"
)

func main() {
	log.Info("starting gateway")

	log.Debug("Loading dotenv")
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	initMinio()
	startExternalApi()
}

func initMinio() {
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))
}

func startExternalApi() {
	fs := service.NewFaktoryService()

	cropHandler := httptransport.NewServer(
		transport.CreateRestHandler(fs, "crop"),
		transport.DecodeCropTask,
		transport.EncodeResponse,
	)
	http.Handle("/crop", cropHandler)

	mostSignificantImageHandler := httptransport.NewServer(
		transport.CreateRestHandler(fs, "most_significant_image"),
		transport.DecodeMostSignificantImageTask,
		transport.EncodeResponse,
	)
	http.Handle("/most_significant_image", mostSignificantImageHandler)

	optimizationHandler := httptransport.NewServer(
		transport.CreateRestHandler(fs, "optimization"),
		transport.DecodeOptimizationTask,
		transport.EncodeResponse,
	)
	http.Handle("/optimization", optimizationHandler)

	portraitHandler := httptransport.NewServer(
		transport.CreateRestHandler(fs, "portrait"),
		transport.DecodePortraitTask,
		transport.EncodeResponse,
	)
	http.Handle("/portrait", portraitHandler)

	screenshotHandler := httptransport.NewServer(
		transport.CreateRestHandler(fs, "screenshot"),
		transport.DecodeScreenshotTask,
		transport.EncodeResponse,
	)
	http.Handle("/screenshot", screenshotHandler)


	fmt.Println(http.ListenAndServe(":8080", nil))
}
