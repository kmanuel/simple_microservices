package main

import (
	"encoding/json"
	"github.com/advancedlogic/GoOse"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

const OutputImageLocation = "../test/output/"

type Request struct {
	In     string `json:"in,omitempty"`
	Out    string `json:"out,omitempty"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", HandleRequest).Methods("POST")
	log.Println(http.ListenAndServe(":8086", router))
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)
	inputUrl := task.In
	outputFile := OutputImageLocation + task.Out

	ExtractMostSignificantImage(inputUrl, outputFile)
}

func ExtractMostSignificantImage(inputUrl string, outputFile string) {
	g := goose.New()
	article, _ := g.ExtractFromURL(inputUrl)
	topImageUrl := article.TopImage
	DownloadImage(topImageUrl, outputFile)
}

func DownloadImage(url string, outputFile string) error {
	filepath := outputFile

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
