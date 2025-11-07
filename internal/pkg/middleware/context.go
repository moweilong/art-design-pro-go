package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	"github.com/moweilong/art-design-pro-go/internal/pkg/contextx"
)

// Context is a middleware that injects common prefix fields to gin.Context.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从当前 span 中获取 traceID 并设置到 gin.Context
		traceID := trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()

		// 将 traceID 存储到新的 context 中，并更新请求的 context
		ctx := contextx.WithTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
