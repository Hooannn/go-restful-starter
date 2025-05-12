package middleware

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func WithTracer(t trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody string
		if body, err := io.ReadAll(c.Request.Body); err == nil {
			requestBody = string(body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		spanName := c.Request.Method + " " + c.FullPath()
		ctx, span := t.Start(c.Request.Context(), spanName)
		defer span.End()

		traceID := span.SpanContext().TraceID().String()
		c.Header("X-Trace-ID", traceID)

		c.Request = c.Request.WithContext(ctx)

		span.SetAttributes(attribute.String("http.request.body", requestBody))

		writer := &responseWriter{ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
		c.Writer = writer

		c.Next()

		span.SetAttributes(attribute.String("http.response.body", writer.body.String()))
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
