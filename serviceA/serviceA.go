package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_a_request_count",
			Help: "Total number of requests to Service A.",
		},
		[]string{"status"},
	)

	latency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "service_a_request_latency_seconds",
			Help:       "Latency of requests to Service A in seconds.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"endpoint"},
	)
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(latency)
}

func main() {
	r := chi.NewRouter()
	r.Use(prometheusMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Post("/cep", GetTemperature)

	// Expose Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Service A running on port 8081...")
	http.ListenAndServe(":8081", r)
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		endpoint := r.URL.Path
		status := strconv.Itoa(w.(interface {
			Status() int
		}).Status())
		latency.WithLabelValues(endpoint).Observe(duration.Seconds())
		requestCount.WithLabelValues(status).Inc()
	})
}

func GetTemperature(w http.ResponseWriter, r *http.Request) {
	cep := strings.TrimSpace(r.FormValue("cep"))

	if len(cep) != 8 || !isNumeric(cep) {
		http.Error(w, "Invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// Simulate processing time
	time.Sleep(200 * time.Millisecond)

	response := struct {
		Cep string `json:"cep"`
	}{
		Cep: cep,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func isNumeric(s string) bool {
	_, err := strconv.ParseUint(s, 10, 64)
	return err == nil
}
