package otel

import (
	"context"
	"encoding/json"
	clientGolangPrometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"go-tracing/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	metric2 "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
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

	RequestCount = clientGolangPrometheus.NewCounterVec(clientGolangPrometheus.CounterOpts{
		Name: "http_request_go_tracing_count",
		Help: "Total number of requset in services go tracing",
	},
		[]string{"url", "status_code"})

	RequestDuration = clientGolangPrometheus.NewHistogramVec(clientGolangPrometheus.HistogramOpts{
		Name:        "http_request_go_tracing_duration",
		Help:        "Duration of request in services go tracing in seconds",
		ConstLabels: nil,
		//Buckets:     clientGolangPrometheus.LinearBuckets(0.001, 0.005, 10),
	},
		[]string{"url", "status_code"})

	ExporterPrometheus   *prometheus.Exporter
	RequestMetricCounter metric2.Int64Counter
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
	clientGolangPrometheus.MustRegister(RequestCount, RequestDuration)
}

func NewPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
}

// InitTracerApp is method to create tracer
func InitTracerApp(ctx context.Context, serviceName string) (*trace.TracerProvider, func()) {
	propagator := NewPropagator()
	otel.SetTextMapPropagator(propagator)

	// create tracer exporter
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

	// create tracer provider
	tracerProvider := NewTraceProvider(exporter, serviceName)
	otel.SetTracerProvider(tracerProvider)

	// assign tracer variable
	OtelApp = &OtelTrace{Trace: tracerProvider.Tracer(serviceName)}

	return tracerProvider, func() {
		_ = tracerProvider.Shutdown(ctx)
	}
}

func (o *OtelTrace) Start(ctx context.Context, name string) (context.Context, otlTrace.Span) {
	return o.Trace.Start(ctx, name)
}

func DumpContext(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	marshal, _ := json.Marshal(&md)
	return string(marshal)
}

func CreatePrometheusExporter() *prometheus.Exporter {
	exporter, err := prometheus.New()
	if err != nil {
		logrus.Fatal(err)
	}

	ExporterPrometheus = exporter
	return ExporterPrometheus
}

func CreatePrometheusMetrixProvider(exp metric.Reader) *metric.MeterProvider {
	provider := metric.NewMeterProvider(metric.WithReader(exp))
	return provider
}

func InitiaizeMetricWithOtelPremetheus(ctx context.Context, serviceName string) func() {
	// create exporter
	exporter := CreatePrometheusExporter()

	// create provider
	provider := CreatePrometheusMetrixProvider(exporter)

	// set global
	otel.SetMeterProvider(provider)

	// get meter
	meter := otel.Meter(serviceName)

	// create counter
	var err error
	RequestMetricCounter, err = meter.Int64Counter("request.total.reo.service", metric2.WithDescription("number of total request http"))
	if err != nil {
		logrus.Fatal(err)
	}

	return func() {
		_ = provider.Shutdown(ctx)
	}
}
