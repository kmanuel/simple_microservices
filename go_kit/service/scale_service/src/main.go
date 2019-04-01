package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/common/log"
	"net/http"
	"os/exec"
	"strconv"
)

type HookMessage struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
}

type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    string            `json:"startsAt,omitempty"`
	EndsAt      string            `json:"EndsAt,omitempty"`
}

var swarmPrefix = "self_impl_swarm_"

type responseJSON struct {
	Status  int
	Message string
}

func asJson(w http.ResponseWriter, status int, message string) {
	data := responseJSON{
		Status:  status,
		Message: message,
	}
	bytes, _ := json.Marshal(data)
	json := string(bytes[:])

	w.WriteHeader(status)
	fmt.Fprint(w, json)
}

func main() {
	myRouter := mux.NewRouter().StrictSlash(false)

	myRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		data := template.Data{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}

		err = handleHookMessage(data)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		asJson(w, http.StatusOK, "success")

	})

	log.Fatal(http.ListenAndServe(":8085", myRouter))
}

func handleHookMessage(message template.Data) error {
	log.Info("received hookMessage123: ", message)

	for _, alert := range message.Alerts {
		alertName := alert.Labels["alertname"]
		log.Info("checking alertName=" + alertName)
		switch alertName {

		case "ManyCropTasksPending":
			scaleUp("crop")
			continue
		case "TooManyCropInstances":
			scaleDown("crop")
			continue
		case "ManyMostSignificantImageTasksPending":
			scaleUp("most_significant_image")
			continue
		case "TooMostSignificantImageInstances":
			return scaleDown("most_significant_image")
			continue
		case "ManyOptimizationTasksPending":
			scaleUp("optimization")
			continue
		case "TooManyOptimizationInstances":
			scaleDown("optimization")
			continue
		case "ManyScreenshotTasksPending":
			scaleUp("screenshot")
			continue
		case "TooManyScreenshotInstances":
			scaleDown("screenshot")
			continue
		}

	}

	return nil
}

func scaleUp(serviceName string) error {
	return scaleTo(serviceName, 5)
}

func scaleDown(serviceName string) error {
	return scaleTo(serviceName, 1)
}

func scaleTo(serviceName string, instanceNum int) error {
	log.Info("scaling service " + serviceName + " to " + strconv.Itoa(instanceNum) + " instances.")
	scaleCmd := exec.Command("docker", "service", "scale", "-d", swarmPrefix+serviceName+"="+strconv.Itoa(instanceNum))
	return scaleCmd.Run()
}
