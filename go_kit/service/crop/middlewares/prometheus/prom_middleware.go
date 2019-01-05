package prometheus

import (
	"flag"
	"fmt"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests handled from faktory.",
		},
		[]string{"controller", "status"},
	)
)

func Init() {
	go start()
}

func start() {
	prometheus.MustRegister(requests)

	var addr = flag.String("listen-address", ":8081", "The address to listen on for HTTP requests.")

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(*addr, nil))
}

type Middleware struct {
	Next    service.CropService
}

func (mw Middleware) CropImage(inputImg string, width int, height int) (output string, err error) {
	requests.With(prometheus.Labels{"controller": "crop", "status": "fetched"}).Inc()
	output, err = mw.Next.CropImage(inputImg, width, height)
	if err != nil {
		requests.With(prometheus.Labels{"controller": "crop", "status": "failed"}).Inc()
	}
	return
}

