package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/gateway/api_image"
	"github.com/prometheus/common/log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	startJsonRestApi()
}

func initMinio() {
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))
}

func startJsonRestApi() {
	myRouter := mux.NewRouter().StrictSlash(false)
	imageHandler := &api_image.ImageHandler{}

	myRouter.HandleFunc("/images", imageHandler.ServeUploadHTTP)
	myRouter.HandleFunc("/images/{id}", imageHandler.ServeDownload)
	myRouter.HandleFunc("/crop", proxyCropRequest)
	myRouter.HandleFunc("/screenshot", proxyScreenshotRequest)
	myRouter.HandleFunc("/most_significant_image", proxyMostSignificantImageRequest)
	myRouter.HandleFunc("/optimization", proxyOptimizationRequest)
	myRouter.HandleFunc("/portrait", proxyPortraitRequest)
	myRouter.HandleFunc("/tasks", proxyRequestServiceRequest)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func proxyRequestServiceRequest(w http.ResponseWriter, r *http.Request) {
	proxyTo("http://request_service:8080", w, r)
}

func proxyScreenshotRequest(w http.ResponseWriter, r *http.Request) {
	proxyTo("http://screenshot:8080", w, r)
}

func proxyPortraitRequest(w http.ResponseWriter, r *http.Request) {
	proxyTo("http://portrait:8080", w, r)
}

func proxyOptimizationRequest(w http.ResponseWriter, r *http.Request) {
	proxyTo("http://optimization:8080", w, r)
}

func proxyMostSignificantImageRequest(w http.ResponseWriter, r *http.Request) {
	proxyTo("http://most_significant_image:8080", w, r)
}

func proxyCropRequest(w http.ResponseWriter, r *http.Request) {
	proxyTo("http://crop:8080", w, r)
}

func proxyTo(proxyTarget string, w http.ResponseWriter, r *http.Request) {
	log.Info("proxying to " + proxyTarget)

	serviceUrl, e := url.Parse(proxyTarget)
	if e != nil {
		panic(e)
	}
	httputil.NewSingleHostReverseProxy(serviceUrl).ServeHTTP(w, r)
}
