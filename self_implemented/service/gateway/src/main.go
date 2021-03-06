package main

import (
	"github.com/gorilla/mux"
	"github.com/kmanuel/simple_microservices/self_implemented/service/gateway/src/api/api_task"
	"github.com/kmanuel/simple_microservices/self_implemented/service/gateway/src/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	dispatchCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dispatch_count",
			Help: "Number of dispatchCounter handled from faktory.",
		},
		[]string{"type"},
	)
)

func main() {
	log.Info("starting gateway")
	go startPrometheus()
	startJsonRestApi()
}

func startPrometheus() {
	log.Debug("starting prometheus")

	prometheus.MustRegister(dispatchCounter)
	registerTaskGauge("crop")
	registerTaskGauge("most_significant_image")
	registerTaskGauge("optimization")
	registerTaskGauge("portrait")
	registerTaskGauge("screenshot")

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func registerTaskGauge(typeName string) {
	faktoryService := service.NewFaktoryService()

	gaugeName := typeName + "_tasks_pending"
	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: gaugeName,
		},
		func() float64 {
			faktoryInfo, err := faktoryService.Info()
			if err != nil {
				log.Error("failed to read FaktoryInfo")
				return 0
			}
			return faktoryInfo.Queues[typeName]
		},
	)); err == nil {
		log.Info("GaugeFunc '" + gaugeName + " registered.")
	}
}

func startJsonRestApi() {
	log.Debug("starting REST API")

	proxyHandler := &api_task.ProxyHandler{DispatchCounter: dispatchCounter}

	myRouter := mux.NewRouter().StrictSlash(false)

	myRouter.HandleFunc("/crop", proxyHandler.CreateCropTask)
	myRouter.HandleFunc("/most_significant_image", proxyHandler.CreateMostSignificantImageTask)
	myRouter.HandleFunc("/optimization", proxyHandler.CreateOptimizationTask)
	myRouter.HandleFunc("/portrait", proxyHandler.CreatePortraitTask)
	myRouter.HandleFunc("/screenshot", proxyHandler.CreateScreenshotTask)
	myRouter.HandleFunc("/info", proxyHandler.GetFaktoryInfo)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
