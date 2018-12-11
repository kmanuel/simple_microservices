package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kmanuel/minioconnector"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

type Request struct {
	In     string `json:"in,omitempty"`
	Out    string `json:"out,omitempty"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Perc   int    `json:"perc"`
	Scale  int    `json:"scale"`
	Blur   int    `json:"blur"`
	Sobel  int    `json:"sobel"`
	Debug  int    `json:"debug"`
	Face   int    `json:"face"`
	Cc     string `json:"cc"`
}

func main() {
	minioconnector.Init()

	router := mux.NewRouter()
	router.HandleFunc("/crop", HandleRequest).Methods("POST")
	log.Println(http.ListenAndServe(":8081", router))
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("received request")
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)

	downloadedFilePath := minioconnector.DownloadFile(task.In)

	outputFilePath := ExtractPortrait(downloadedFilePath, task.Width, task.Height)

	minioconnector.UploadFile(outputFilePath)

	log.Printf("Command finished")
}

func ExtractPortrait(
	inputLocation string,
	width int,
	height int) string {

	outputFilePath := "/tmp/" + uuid.New().String() + ".jpg"

	cmd := exec.Command(
		"caire",
		"-in", inputLocation,
		"-out", outputFilePath,
		"-width="+strconv.Itoa(width),
		"-height="+strconv.Itoa(height),
		"-perc=1",
		"-square=0",
		"-scale=1",
		"-blur=0",
		"-sobel=0",
		"-debug=0",
		"-face=1",
		"-cc=./data/facefinder",
	)
	cmd.Run()

	return outputFilePath
}
