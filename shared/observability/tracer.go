package observability

//What tracer.go Actually Is
//tracer.go is not a template. It is not magic. It is just a Go file that configures the OpenTelemetry SDK and tells it two things:
//Copy1. Where to send traces (Alloy on port 4318)
//2. What to call this service (city-explorer-api)

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// InitTracer sets up the OpenTelemetry tracer.
// It returns a shutdown function you must call
// when the service exits to flush pending traces.
//
// Call it like this in main():
//
//	shutdown, err := telemetry.InitTracer(ctx, "api")
//	if err != nil { log.Fatal(err) }
//	defer shutdown(ctx)
func InitTracer(ctx context.Context, serviceName string) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT not set")
	}

	// The exporter sends traces to Alloy
	// Alloy forwards them to Grafana Cloud Tempo
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithHeaders(map[string]string{
			"Content-Type": "application/json",
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Resource describes this service to Grafana
	// These labels appear on every trace
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("0.1.0"),
			semconv.ServiceNamespace("city-explorer"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// TracerProvider is the core of the tracing setup
	// BatchSpanProcessor sends spans in batches
	// rather than one at a time — more efficient
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	//Add Propagator To tracer.go
	//Since both services call observability.InitTracer() this registers the propagator for both services automatically. You do not need to add it to each main.go separately.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Register as the global tracer provider
	// Now any code can call otel.Tracer() to get a tracer
	otel.SetTracerProvider(tp)
	slog.Info("tracer initialized")

	return tp.Shutdown, nil
}
