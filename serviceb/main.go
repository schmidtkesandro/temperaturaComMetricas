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
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	otlptracegrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.8.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var tracer trace.Tracer

func initTracer(serviceName string) func() {
	otel.SetTextMapPropagator(propagation.TraceContext{})

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
	shutdownB := initTracer("serviceB")
	defer shutdownB()

	r := mux.NewRouter()
	r.Use(recordMetrics)

	r.Handle("/cep", otelhttp.NewHandler(http.HandlerFunc(handleCEP), "ReceiveCEP"))
	r.Handle("/metrics", promhttp.Handler())

	log.Println("Service B is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func handleCEP(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Received Headers:", r.Header)

	ctx, span := tracer.Start(r.Context(), "HandleCEP")
	defer span.End()

	var req struct {
		CEP string `json:"cep"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	location, err := getLocation(ctx, req.CEP)

	if err != nil {
		http.Error(w, "CEP not found: "+err.Error(), http.StatusNotFound)
		return
	}

	temperature, err := getTemperature(ctx, location)
	if err != nil {
		http.Error(w, "Failed: "+err.Error(), http.StatusInternalServerError)
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
	_, spanGetLocation := tracer.Start(ctx, "getLocation")
	defer spanGetLocation.End()

	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))

	if err != nil {
		return "", fmt.Errorf("CEP not found")
	}
	defer resp.Body.Close()
	fmt.Println("Status de Retorno CEP", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get location")
	}

	var data struct {
		Localidade string `json:"localidade"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	fmt.Println("err ", err, "localidade ", data.Localidade)
	if err != nil || data.Localidade == "" {
		return "", fmt.Errorf("CEP not found")
	}

	return data.Localidade, nil
}

func getTemperature(ctx context.Context, city string) (float64, error) {
	_, spanGetTemperature := tracer.Start(ctx, "getTemperature")
	defer spanGetTemperature.End()

	apiKey := os.Getenv("WEATHER_API_KEY")
	escapedCity := url.QueryEscape(city)
	fmt.Println(escapedCity)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, escapedCity)

	resp, err := http.Get(url)
	fmt.Println("Status de Retorno temperatura", resp.StatusCode)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get temperature - status %d", resp.StatusCode)
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
