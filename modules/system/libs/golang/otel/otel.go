package otel

import (
	"context"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

type telemetry struct {
	client otlptrace.Client
	tp     *trace.TracerProvider
}

type Telemetry interface {

	// Gets used trace provider
	// GetTraceProvider() *trace.TracerProvider
	// Sends buffered data and closes connections
	Shutdown(ctx context.Context)
}

/*
	Registers OpenTelemetry and adds global hooks.
*/
func Register(ctx context.Context, endpoint string, serviceModule string, serviceName string, serviceVersion string, serviceInstanceID string) (Telemetry, error) {
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	exp, err := otlptrace.New(ctx, client)

	if err != nil {
		return nil, err
	}

	r, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
			semconv.ServiceNamespaceKey.String(serviceModule),
			semconv.ServiceInstanceIDKey.String(serviceInstanceID),
		),
	)

	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(r),
		trace.WithSampler(trace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return &telemetry{client: client, tp: tp}, nil
}

func (t *telemetry) Shutdown(ctx context.Context) {
	t.tp.Shutdown(ctx)
	t.client.Stop(ctx)
}
