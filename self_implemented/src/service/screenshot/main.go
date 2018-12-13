package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
)

type Request struct {
	Url string `json:"url"`
}

func main() {
	godotenv.Load()
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))

	router := mux.NewRouter()
	router.HandleFunc("/", HandleRequest).Methods("POST")
	log.Info(http.ListenAndServe(":8080", router))
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	log.Info("received request")
	json.NewEncoder(w).Encode("received request")
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)

	outputFilePath := takeScreenShot(task.Url)

	minioconnector.UploadFile(outputFilePath)

	log.Info("finished request")
}

func takeScreenShot(url string) string {
	log.WithField("url", url).Info("taking screenshot")

	chromeUserAgent := "Mozilla/5.0 (Windows NT 6.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	phantomJSBin := "./lib/bin/phantomjs"
	jsPath := "./lib/js/screenshot.js"
	logFile := "screenshot_service.log"

	outputFilePath := "/tmp/output" + uuid.New().String() + ".jpg"

	cmd := exec.Command(phantomJSBin, jsPath, url, outputFilePath, logFile, chromeUserAgent)

	if err := cmd.Run(); nil != err {
		log.Printf("process job err - %s\n", err.Error())
	}

	log.WithField("url", url).Info("took screenshot")

	return outputFilePath
}
