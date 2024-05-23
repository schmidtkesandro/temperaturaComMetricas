package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

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
	shutdown := initTracer("serviceA")
	defer shutdown()

	r := mux.NewRouter()
	r.Use(recordMetrics)

	r.Handle("/", otelhttp.NewHandler(http.HandlerFunc(homeHandler), "Index"))
	r.Handle("/cep", otelhttp.NewHandler(http.HandlerFunc(handleCEP), "GetCEP"))
	r.Handle("/metrics", promhttp.Handler())

	log.Println("Service A is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleCEP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "HandleCEP")
	defer span.End()

	var req struct {
		CEP string `json:"cep"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if !isValidCEP(req.CEP) {
		http.Error(w, "invalid CEP format; must be 8 numeric characters", http.StatusBadRequest)
		return
	}
	cepData := map[string]string{"cep": req.CEP}
	cepBytes, err := json.Marshal(cepData)

	if err != nil {
		http.Error(w, "failed to marshal CEP data", http.StatusInternalServerError)
		return
	}
	resp, err := makeRequest(ctx, "http://serviceb:8081/cep", cepBytes)
	if err != nil {
		http.Error(w, "Failed to call service B: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Status de Retorno", resp.StatusCode)
	// if resp.StatusCode != http.StatusOK {
	// 	http.Error(w, "service B error", resp.StatusCode)
	// 	return
	// }
	// Read the response body and status code
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from service B", http.StatusInternalServerError)
		return
	}

	// Write the status code and response body from service B to the client
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
	// var responseBody map[string]interface{}
	// if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
	// 	http.Error(w, "failed to read response from service B", http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(responseBody)
}
func isValidCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}
func makeRequest(ctx context.Context, url string, data []byte) (*http.Response, error) {
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Explicitly ensure context is propagated
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
	//fmt.Println("Sending Headers:", req.Header)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
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
