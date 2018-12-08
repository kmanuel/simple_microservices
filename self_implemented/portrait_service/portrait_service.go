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
	Square int    `json:"square"`
	Scale  int    `json:"scale"`
	Blur   int    `json:"blur"`
	Sobel  int    `json:"sobel"`
	Debug  int    `json:"debug"`
	Face   int    `json:"face"`
	Cc     string `json:"cc"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/crop", ExtractPortrait).Methods("POST")
	log.Println(http.ListenAndServe(":8081", router))
}

func ExtractPortrait(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("received request")
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)

	log.Println(task)

	cmd := exec.Command(
		"caire",
		"-in", InputImageLocation+task.In,
		"-out", OutputImageLocation+task.Out,
		"-width="+strconv.Itoa(task.Width),
		"-height="+strconv.Itoa(task.Height),
		"-perc=1",
		"-square="+strconv.Itoa(task.Square),
		"-scale="+strconv.Itoa(task.Scale),
		"-blur="+strconv.Itoa(task.Blur),
		"-sobel="+strconv.Itoa(task.Sobel),
		"-debug="+strconv.Itoa(task.Debug),
		"-face="+strconv.Itoa(task.Face),
		"-cc=./data/facefinder",
	)
	log.Println(cmd)
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
}