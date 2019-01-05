package main

import (
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/kmanuel/simple_microservices/go_kit/service/request_service/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/request_service/transport"
	"net/http"
)

func main() {

	var ss service.RequestStatusService
	ss = service.RequestStatusServiceImpl{}

	requestListHandler := httptransport.NewServer(
		transport.MakeStatusEndpoint(ss),
		transport.DecodeListRequest,
		transport.EncodeResponse,
	)

	http.Handle("/", requestListHandler)
	fmt.Println(http.ListenAndServe(":8080", nil))

}
