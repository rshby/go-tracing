package test

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"testing"
)

func TestOtelPrometheus(t *testing.T) {
	// Create a new Prometheus exporter
	exporter, err := prometheus.New()
	if err != nil {
		panic(err)
	}

	// Create a meter provider and register the Prometheus exporter
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	defer provider.Shutdown(context.Background())

	// Set the global meter provider
	otel.SetMeterProvider(provider)

	// Initialize your Gin engine
	app := gin.Default()

	// Expose the /metrics endpoint using the Prometheus exporter handler
	app.GET("/metrics", gin.WrapH(exporter))

	// Run the Gin server
	app.Run(":8080")
}
