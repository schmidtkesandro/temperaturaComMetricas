package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

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

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shut down tracer provider: %v", err)
		}
	}
}
func main() {
	// Service B
	shutdown := initTracer("serviceB")
	defer shutdown()
	r := mux.NewRouter()
	r.Use(recordMetrics)
	r.HandleFunc("/cep", handleCEP).Methods("POST")
	r.Handle("/metrics", promhttp.Handler())

	log.Println("Service B is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
