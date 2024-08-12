package otel

import (
	"context"
	"encoding/json"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"go-tracing/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	// "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
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

	RequestCount = prometheus2.NewCounter(prometheus2.CounterOpts{
		Name: "http_request_go_tracing_count",
		Help: "Total number of requset in services go tracing",
	})

	RequestDuration = prometheus2.NewHistogram(prometheus2.HistogramOpts{
		Name:        "http_request_go_tracing_duration",
		Help:        "Duration of request in services go tracing in seconds",
		ConstLabels: nil,
		Buckets:     prometheus2.LinearBuckets(0.001, 0.005, 10),
	})
)

// NewTraceExporter is method to create exporter jaeger
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

// NewConsoleExporter is method to create exporter console
func NewConsoleExporter(ctx context.Context) (trace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

func NewTraceProvider(exporter trace.SpanExporter, serviceName string) *trace.TracerProvider {
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(1*time.Second)),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(serviceName))),
		trace.WithSampler(trace.AlwaysSample()))

	return traceProvider
}

// NewMetrixPrometheus is function to create and set metrics prometheus
func NewMetrixPrometheus(ctx context.Context, name string) {
	prometheus2.MustRegister(RequestCount, RequestDuration)
}

// InitTracerApp is method to
func InitTracerApp(ctx context.Context, serviceName string) (*trace.TracerProvider, func()) {
	var exporter trace.SpanExporter
	switch config.OtelExporter() {
	case "console":
		var err error
		exporter, err = NewConsoleExporter(ctx)
		if err != nil {
			logrus.Fatalf("failed to get console exporter : %s", err.Error())
		}
	case "jaeger":
		var err error
		exporter, err = NewTraceExporter(ctx)
		if err != nil {
			logrus.Fatalf("failed to get exporter : %s", err.Error())
		}
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
