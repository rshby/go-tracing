package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type customWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func NewCustomeWritter(r gin.ResponseWriter) *customWriter {
	return &customWriter{
		ResponseWriter: r,
		body:           bytes.NewBufferString(""),
	}
}

func (c *customWriter) Writer(b []byte) (int, error) {
	c.body.Write(b)
	return c.ResponseWriter.Write(b)
}
