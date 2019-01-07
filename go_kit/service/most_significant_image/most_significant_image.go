package main

import (
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/joho/godotenv"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/transport"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	s := service.NewMostSignificantImageService()
	startRestApi(&s)
}

func startRestApi(s *service.MostSignificantImageService) {
	requestHandler := httptransport.NewServer(
		transport.MakeMostSignificantImageEndpoint(*s),
		transport.DecodeMostSignificantImageTask,
		transport.EncodeResponse,
	)
	http.Handle("/", requestHandler)
	fmt.Println(http.ListenAndServe(":"+os.Getenv("MOST_SIGNIFICANT_IMAGE_EXTERNAL_PORT"), nil))
}
