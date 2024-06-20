package main

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"log"
	"os"
)

func main() {
	hostname, _ := os.Hostname()
	ctx := context.Background()

	ingestKey := os.Getenv("TELEMETRYHUB_KEY")

	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("myService"),
		semconv.ServiceVersion("0.0.1"),
		semconv.HostName(hostname),
		attribute.String("telemetryhub.key", ingestKey),
	)

	exporter, _ := otlptracegrpc.New(
		// Use an appropriate context here, e.g. context.Background() since this is probably
		// being constructed at boot time.
		ctx,
		// Use with local collector
		// otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("otlp.telemetryhub.com:4317"),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
		//otlptracegrpc.WithHeaders(headers),
	)

	// Create the tracer provider with the exporter and SDK options
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		// Resource constructed earlier
		sdktrace.WithResource(r),
	)
	defer func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shutdown the tracer provider: %v", err)
		}
	}()
}
