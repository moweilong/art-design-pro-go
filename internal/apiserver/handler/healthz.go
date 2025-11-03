package handler

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/onexstack/onexstack/pkg/core"

	v1 "github.com/moweilong/art-design-pro-go/pkg/api/apiserver/v1"
)

// Healthz 服务健康检查.
func (h *Handler) Healthz(c *gin.Context) {
	slog.InfoContext(c.Request.Context(), "Healthz handler is called", "method", "Healthz", "status", "healthy")
	core.WriteResponse(c, v1.HealthzResponse{
		Status:    v1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil)
}
