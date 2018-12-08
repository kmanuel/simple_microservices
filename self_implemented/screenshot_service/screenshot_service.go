package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/exec"
)

type Request struct {
	Url string `json:"url"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", HandleRequest).Methods("POST")
	log.Println(http.ListenAndServe(":8083", router))
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("received request")
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)

	takeScreenShot(task.Url)
}

func takeScreenShot(url string) {
	chromeUserAgent := "Mozilla/5.0 (Windows NT 6.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	phantomJSBin := "./lib/bin/phantomjs"
	jsPath := "./lib/js/screenshot_service.js"
	logFile := "screenshot_service.log"

	cmd := exec.Command(phantomJSBin, jsPath, url, "/app/output/testoutput.jpg", logFile, chromeUserAgent)

	if err := cmd.Run(); nil != err {
		log.Printf("process job err - %s\n", err.Error())
	}
}
