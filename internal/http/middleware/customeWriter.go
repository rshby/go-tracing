package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type customWriter struct {
	StatusCode int
	gin.ResponseWriter
	Body *bytes.Buffer
}

func NewCustomeWritter(r gin.ResponseWriter) *customWriter {
	return &customWriter{
		ResponseWriter: r,
		Body:           bytes.NewBufferString(""),
	}
}

func (c *customWriter) Write(b []byte) (int, error) {
	c.Body.Write(b)
	return c.ResponseWriter.Write(b)
}

func (c *customWriter) WriteHeader(statusCode int) {
	c.StatusCode = statusCode
	c.ResponseWriter.WriteHeader(statusCode)
}
