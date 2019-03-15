package main

import (
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/kmanuel/simple_microservices/go_kit/service/gateway/src/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/gateway/src/transport"
	"github.com/prometheus/common/log"
	"net/http"
)

func main() {
	log.Info("starting gateway")

	startExternalApi()
}

func startExternalApi() {
	fs := service.NewFaktoryService()

	cropHandler := httptransport.NewServer(
		transport.CreateFaktoryHandler(fs, "crop"),
		transport.DecodeCropTask,
		transport.EncodeResponse,
	)
	http.Handle("/crop", cropHandler)

	mostSignificantImageHandler := httptransport.NewServer(
		transport.CreateFaktoryHandler(fs, "most_significant_image"),
		transport.DecodeMostSignificantImageTask,
		transport.EncodeResponse,
	)
	http.Handle("/most_significant_image", mostSignificantImageHandler)

	optimizationHandler := httptransport.NewServer(
		transport.CreateFaktoryHandler(fs, "optimization"),
		transport.DecodeOptimizationTask,
		transport.EncodeResponse,
	)
	http.Handle("/optimization", optimizationHandler)

	portraitHandler := httptransport.NewServer(
		transport.CreateFaktoryHandler(fs, "portrait"),
		transport.DecodePortraitTask,
		transport.EncodeResponse,
	)
	http.Handle("/portrait", portraitHandler)

	screenshotHandler := httptransport.NewServer(
		transport.CreateFaktoryHandler(fs, "screenshot"),
		transport.DecodeScreenshotTask,
		transport.EncodeResponse,
	)
	http.Handle("/screenshot", screenshotHandler)

	infoHandler := httptransport.NewServer(
		transport.CreateRestHandler(fs),
		transport.DecodeInfoRequest,
		transport.EncodeResponse,
	)
	http.Handle("/info", infoHandler)

	fmt.Println(http.ListenAndServe(":8080", nil))
}
