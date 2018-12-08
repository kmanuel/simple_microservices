package main

import (
"encoding/json"
"fmt"
"github.com/gorilla/mux"
"github.com/muesli/smartcrop"
"github.com/muesli/smartcrop/nfnt"
"image"
"image/jpeg"
"log"
"net/http"
"os"
)

const InputImageLocation = "../test/input/"
const OutputImageLocation = "../test/output/"

type Request struct {
	In string `json:"in,omitempty"`
	Out string `json:"out,omitempty"`
	Width int `json:"width"`
	Height int `json:"height"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", HandleRequest).Methods("POST")
	log.Println(http.ListenAndServe(":8082", router))
}

func HandleRequest(w http.ResponseWriter, r * http.Request) {
	json.NewEncoder(w).Encode("received request")
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)

	inputFile := task.In
	outputFile := task.Out

	optimizeImage(inputFile, outputFile)

	log.Printf("finished")
}

func optimizeImage(inputFile string, outputFile string) {
	f, _ := os.Open(inputFile)
	img, _, _ := image.Decode(f)
	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	topCrop, _ := analyzer.FindBestCrop(img, 250, 250)
	// The crop will have the requested aspect ratio, but you need to copy/scale it yourself
	fmt.Printf("Top crop: %+v\n", topCrop)
	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	croppedImg := img.(SubImager).SubImage(topCrop)
	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, croppedImg, nil)
}

