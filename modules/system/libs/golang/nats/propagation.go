package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type headersPropagator struct {
	Msg *nats.Msg
}

func (c *headersPropagator) Get(key string) string {
	return c.Msg.Header.Get(key)
}
func (c *headersPropagator) Set(key string, value string) {
	c.Msg.Header.Set(key, value)
}
func (c *headersPropagator) Keys() []string {
	keys := make([]string, 0, len(c.Msg.Header))
	for k := range c.Msg.Header {
		keys = append(keys, k)
	}
	return keys
}

func InjectTelemetryContext(ctx context.Context, msg *nats.Msg) {
	propagator := headersPropagator{Msg: msg}
	otel.GetTextMapPropagator().Inject(ctx, &propagator)
}

func RetrieveTelemetryContext(ctx context.Context, msg *nats.Msg) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, &headersPropagator{Msg: msg})
}

func StartTelemetrySpanFromMessage(ctx context.Context, msg *nats.Msg, spanName string) (context.Context, trace.Span) {
	telemetryContext := RetrieveTelemetryContext(ctx, msg)
	return otel.GetTracerProvider().Tracer("system_nats/otellib").Start(telemetryContext, spanName)
}
