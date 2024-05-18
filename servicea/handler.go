package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(requestDuration)
}

func recordMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(requestDuration.WithLabelValues(r.URL.Path))
		defer timer.ObserveDuration()
		next.ServeHTTP(w, r)
	})
}

type CEPRequest struct {
	CEP string `json:"cep"`
}

func handleCEP(w http.ResponseWriter, r *http.Request) {
	var req CEPRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil || len(req.CEP) != 8 {
		http.Error(w, `{"error":"invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	resp, err := http.Post("http://serviceb:8081/cep", "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, `{"error":"failed to connect to service B"}`, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(respBody))
		http.Error(w, `{"error":"service B error"}`, resp.StatusCode)
		return
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to read response from service B"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}
