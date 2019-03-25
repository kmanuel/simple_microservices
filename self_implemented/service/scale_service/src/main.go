package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
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

func main() {
	myRouter := mux.NewRouter().StrictSlash(false)

	myRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var hookMessage HookMessage
		err := dec.Decode(&hookMessage)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}

		err = handleHookMessage(hookMessage)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	log.Fatal(http.ListenAndServe(":8085", myRouter))
}

func handleHookMessage(message HookMessage) error {
	log.Info("received hookMessage: ", message)

	for _, alert := range message.Alerts {
		alertName := alert.Labels["alertname"]
		log.Info("checking alertName=" + alertName)
		switch alertName {
		case "ManyCropTasksPending":
			return scaleUp("crop")
		case "ManyMostSignificantImageTasksPending":
			return scaleUp("most_significant_image")
		case "ManyOptimizationTasksPending":
			return scaleUp("optimization")
		case "ManyScreenshotTasksPending":
			return scaleUp("screenshot")
		case "TooManyCropInstances":
			return scaleDown("crop")
		case "TooManyMostSignificantImageInstances":
			return scaleDown("most_significant_image")
		case "TooManyOptimizationInstances":
			return scaleDown("optimization")
		case "TooManyScreenshotInstances":
			return scaleDown("screenshot")
		}
	}

	log.Info("found no matching service to scale")
	return nil
}

func scaleUp(serviceName string) error {
	return scaleTo(serviceName, 5)
}

func scaleDown(serviceName string) error {
	return scaleTo(serviceName, 1)
}

func scaleTo(serviceName string, instanceNum int) error {
	scaleCmd := exec.Command("docker", "service", "scale", "-d", "self_impl_swarm_"+serviceName+"="+strconv.Itoa(instanceNum))
	return scaleCmd.Run()
}
