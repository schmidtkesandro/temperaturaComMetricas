package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.8.0"
)

func initTracer(serviceName string) func() {
	endpoint := "http://localhost:9411/api/v2/spans"

	exporter, err := zipkin.New(
		endpoint,
	)
	if err != nil {
		log.Fatalf("failed to initialize Zipkin exporter %v", err)
	}

	bsp := trace.NewBatchSpanProcessor(exporter)
	tracerProvider := trace.NewTracerProvider(
		trace.WithSpanProcessor(bsp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shut down tracer provider %v", err)
		}
	}
}

func main() {
	// Service A
	shutdown := initTracer("serviceA")
	defer shutdown()
	r := mux.NewRouter()
	r.HandleFunc("/cep", handleCEP).Methods("POST")
	r.HandleFunc("/home", homeHandler)

	log.Println("Service A is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
