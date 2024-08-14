package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-tracing/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		startTime := time.Now()

		ctx := c.Request.Context()

		// Check if the context already contains a span
		span := trace.SpanFromContext(ctx)
		if span.SpanContext().IsValid() {
			// If there is already a valid span, proceed with it
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			return
		}

		// Start a new trace if no valid span is found
		url := fmt.Sprintf("[%s] %s", c.Request.Method, c.Request.RequestURI)
		ctx, span = otel.OtelApp.Start(ctx, url)

		// Add attributes to the span
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.RequestURI),
			attribute.String("http.client_ip", c.ClientIP()),
		)

		// End the span when the request is done
		defer span.End()

		// Replace the context in the request with the one containing the span
		c.Request = c.Request.WithContext(ctx)

		// catch response
		customWriters := NewCustomeWritter(c.Writer)
		c.Writer = customWriters

		// Proceed with the request
		c.Next()

		// Set additional attributes based on the response
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)

		// for metrics
		otel.RequestCount.WithLabelValues(c.Request.RequestURI, fmt.Sprintf("%d", customWriters.Status())).Inc()
		elapseTime := time.Since(startTime)
		otel.RequestDuration.WithLabelValues(c.Request.RequestURI, fmt.Sprintf("%d", customWriters.Status())).Observe(elapseTime.Seconds())
	}
}
