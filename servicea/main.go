package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

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
	shutdown := initTracer("serviceA")
	defer shutdown()

	r := mux.NewRouter()
	r.Use(recordMetrics)
	r.HandleFunc("/cep", handleCEP).Methods("POST")
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/", homeHandler)

	log.Println("Service A is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
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

	// Inicia um novo span para a chamada ao serviço B
	_, spanCallServiceB := tracer.Start(ctx, "callServiceB")
	defer spanCallServiceB.End()

	// Prepare the request body for the service B
	requestBody, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "failed to create request body", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://serviceb:8081/cep", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "failed to call service B", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "service B error", resp.StatusCode)
		return
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		http.Error(w, "failed to read response from service B", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseBody)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func recordMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record metrics here if necessary
		next.ServeHTTP(w, r)
	})
}
