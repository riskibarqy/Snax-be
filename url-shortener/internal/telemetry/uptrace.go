package telemetry

import (
	"context"
	"log"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// InitUptrace initializes Uptrace with the provided DSN
func InitUptrace(dsn string) error {
	// Configure Uptrace
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(dsn),
		uptrace.WithServiceName("url-shortener"),
		uptrace.WithServiceVersion("1.0.0"),
	)

	// Set up a tracer
	tracer = otel.Tracer("url-shortener")

	return nil
}

// Shutdown gracefully shuts down OpenTelemetry
func Shutdown(ctx context.Context) {
	if err := uptrace.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down Uptrace: %v", err)
	}
}

// StartSpan starts a new span with the given name and returns the context and span
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}

// AddSpanError adds an error to the current span
func AddSpanError(span trace.Span, err error) {
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		span.RecordError(err)
	}
}

// GetTracer returns the global tracer
func GetTracer() trace.Tracer {
	return tracer
}
