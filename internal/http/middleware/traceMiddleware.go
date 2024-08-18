package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-tracing/otel"
	"go.opentelemetry.io/otel/attribute"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		url := fmt.Sprintf("[%s] %s", c.Request.Method, c.Request.RequestURI)
		ctx, span := otel.OtelApp.Start(ctx, url)

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
		customWriters := &customWriter{
			ResponseWriter: c.Writer,
			Body:           &bytes.Buffer{},
		}
		c.Writer = customWriters

		// Proceed with the request
		c.Next()

		// Set additional attributes based on the response
		span.SetAttributes(
			attribute.Int("http.status_code", customWriters.StatusCode),
		)

		if customWriters.Body.String() != "" {
			span.SetAttributes(attribute.String("http.response_body", customWriters.Body.String()))
		}
	}
}
