package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-tracing/database"
	"go-tracing/internal/config"
	"go-tracing/internal/http/router"
	"go-tracing/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	trace2 "go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	t trace2.Tracer
)

func init() {
	logger.SetupLogger()
}

func main() {
	mysqlDB, mysqlCloser := database.InitializeMysqlDatabase()
	defer mysqlCloser()
	logrus.Info(mysqlDB)

	traceProvider, shutdownTrace := initTracerApp(context.Background(), "go-traciing")
	defer shutdownTrace()

	t = traceProvider.Tracer("go-tracing")
	app := gin.Default()

	// router
	router.NewRouter(&app.RouterGroup)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.AppPort()),
		Handler: app,
	}

	var (
		wg         = &sync.WaitGroup{}
		chanSignal = make(chan os.Signal)
	)

	signal.Notify(chanSignal, os.Interrupt)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			select {
			case <-chanSignal:
				_ = server.Close()
				return
			}
		}
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		logrus.Infof("running on port %d", config.AppPort())
		if err := server.ListenAndServe(); err != nil {
			logrus.Error(err)
			return
		}
	}(wg)

	wg.Wait()
}

func newTraceExporter(ctx context.Context) (trace.SpanExporter, error) {
	exporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithHeaders(map[string]string{
			"content-type": "application/json",
		}),
		otlptracehttp.WithInsecure()))

	return exporter, err
}

func newTraceProvider(exporter trace.SpanExporter, serviceName string) *trace.TracerProvider {
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(1*time.Second)),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(serviceName))))

	return traceProvider
}

func initTracerApp(ctx context.Context, serviceName string) (*trace.TracerProvider, func()) {
	exporter, _ := newTraceExporter(ctx)

	tracerProvideer := newTraceProvider(exporter, serviceName)
	otel.SetTracerProvider(tracerProvideer)

	return tracerProvideer, func() {
		_ = tracerProvideer.Shutdown(ctx)
	}
}

func New(ctx context.Context, name string) (context.Context, trace2.Span) {
	return t.Start(ctx, name)
}

func NewSpan() {

}
