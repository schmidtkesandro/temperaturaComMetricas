package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	otlptracegrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.8.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var tracer trace.Tracer

func initTracer(serviceName string) func() {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "otelcol:4317", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create gRPC connection to collector: %v", err)
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tracerProvider)
	tracer = otel.Tracer(serviceName)

	return func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shut down tracer provider: %v", err)
		}
	}
}

func main() {
	shutdown := initTracer("serviceB")
	defer shutdown()

	r := mux.NewRouter()
	r.Use(recordMetrics)
	r.HandleFunc("/cep", handleCEP).Methods("POST")
	r.Handle("/metrics", promhttp.Handler())

	log.Println("Service B is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func handleCEP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleCEP")
	defer span.End()

	var req struct {
		CEP string `json:"cep"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Inicia um novo span para a busca da localização
	ctx, spanLocation := tracer.Start(ctx, "getLocation")
	location, err := getLocation(ctx, req.CEP)
	spanLocation.End()
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	// Inicia um novo span para a busca da temperatura
	ctx, spanTemperature := tracer.Start(ctx, "getTemperature")
	temperature, err := getTemperature(ctx, location)
	spanTemperature.End()
	if err != nil {
		http.Error(w, "failed to get temperature", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"city":   location,
		"temp_C": temperature,
		"temp_F": temperature*1.8 + 32,
		"temp_K": temperature + 273,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getLocation(ctx context.Context, cep string) (string, error) {
	_, span := tracer.Start(ctx, "getLocation")
	defer span.End()

	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get location")
	}

	var data struct {
		Localidade string `json:"localidade"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data.Localidade, nil
}

func getTemperature(ctx context.Context, city string) (float64, error) {
	_, span := tracer.Start(ctx, "getTemperature")
	defer span.End()

	apiKey := os.Getenv("WEATHER_API_KEY")
	escapedCity := url.QueryEscape(city)
	fmt.Println(escapedCity)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, escapedCity)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get temperature")
	}

	var weatherResponse struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		return 0, err
	}

	return weatherResponse.Current.TempC, nil
}

func recordMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record metrics here if necessary
		next.ServeHTTP(w, r)
	})
}
