package otel

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	otlTrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"time"
)

type OtelTrace struct {
	Trace otlTrace.Tracer
}

var (
	OtelApp = &OtelTrace{}
)

func NewTraceExporter(ctx context.Context) (trace.SpanExporter, error) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
		otlptracehttp.WithHeaders(map[string]string{
			"Content-Type": "application/json",
		}),
		otlptracehttp.WithInsecure())

	return exporter, err
}

func NewTraceProvider(exporter trace.SpanExporter, serviceName string) *trace.TracerProvider {
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(1*time.Second)),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(serviceName))),
		trace.WithSampler(trace.AlwaysSample()))

	return traceProvider
}

func InitTracerApp(ctx context.Context, serviceName string) (*trace.TracerProvider, func()) {
	exporter, err := NewTraceExporter(ctx)
	if err != nil {
		logrus.Fatalf("failed to get exporter : %s", err.Error())
	}

	tracerProvideer := NewTraceProvider(exporter, serviceName)
	otel.SetTracerProvider(tracerProvideer)

	OtelApp = &OtelTrace{Trace: tracerProvideer.Tracer(serviceName)}

	return tracerProvideer, func() {
		_ = tracerProvideer.Shutdown(ctx)
	}
}

func (o *OtelTrace) Start(ctx context.Context, name string) (context.Context, otlTrace.Span) {
	return o.Trace.Start(ctx, name, otlTrace.WithAttributes(attribute.String("context", DumpContext(ctx))))
}

func DumpContext(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	marshal, _ := json.Marshal(&md)
	return string(marshal)
}
