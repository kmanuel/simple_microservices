package api_task

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyHandler struct {
	DispatchCounter *prometheus.CounterVec
}

func (h *ProxyHandler) ProxyToRequestService(w http.ResponseWriter, r *http.Request) {
	log.Info("received request for all tasks")
	h.DispatchCounter.With(prometheus.Labels{"type": "request_service"}).Inc()
	proxyTo("http://request_service:8080", w, r)
}

func (h *ProxyHandler) ProxyToCropService(w http.ResponseWriter, r *http.Request) {
	log.Info("proxying to crop service")
	h.DispatchCounter.With(prometheus.Labels{"type": "crop"}).Inc()
	proxyTo("http://crop:8080", w, r)
}

func (h *ProxyHandler) ProxyToScreenshotService(w http.ResponseWriter, r *http.Request) {
	log.Info("proxying to screenshot service")
	h.DispatchCounter.With(prometheus.Labels{"type": "screenshot"}).Inc()
	proxyTo("http://screenshot:8080", w, r)
}

func (h *ProxyHandler) ProxyToMostSignificantImageService(w http.ResponseWriter, r *http.Request) {
	log.Info("proxying to most_significant_image service")
	h.DispatchCounter.With(prometheus.Labels{"type": "most_significant_image"}).Inc()
	proxyTo("http://most_significant_image:8080", w, r)
}

func (h *ProxyHandler) ProxyToOptimizationService(w http.ResponseWriter, r *http.Request) {
	log.Info("proxying to optimization service")
	h.DispatchCounter.With(prometheus.Labels{"type": "optimization"}).Inc()
	proxyTo("http://optimization:8080", w, r)
}

func (h *ProxyHandler) ProxyToPortraitService(w http.ResponseWriter, r *http.Request) {
	log.Info("proxying to portrait service")
	h.DispatchCounter.With(prometheus.Labels{"type": "portrait"}).Inc()
	proxyTo("http://portrait:8080", w, r)
}

func proxyTo(proxyTargetHost string, w http.ResponseWriter, r *http.Request) {
	requestServiceUrl, e := url.Parse(proxyTargetHost)
	if e != nil {
		panic(e)
	}
	httputil.NewSingleHostReverseProxy(requestServiceUrl).ServeHTTP(w, r)
}
