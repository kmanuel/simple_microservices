package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

const InputImageLocation = "../test/input/"
const OutputImageLocation = "../test/output/"

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
	router := mux.NewRouter()
	router.HandleFunc("/crop", HandleRequest).Methods("POST")
	log.Println(http.ListenAndServe(":8081", router))
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("received request")
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)

	inputLocation := InputImageLocation + task.In
	outputLocation := OutputImageLocation + task.Out

	err := ExtractPortrait(inputLocation, outputLocation, task.Width, task.Height)

	log.Printf("Command finished with error: %v", err)
}

func ExtractPortrait(
	inputLocation string,
	outputLocation string,
	width int,
	height int) error {

	cmd := exec.Command(
		"caire",
		"-in", inputLocation,
		"-out", outputLocation,
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
	log.Println(cmd)
	err := cmd.Run()
	return err
}
