package httpserver

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"fragments/internal/core"
)

func LogAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		id, _ := core.NewId() //uid.NewMono()
		c.Set("request-id", id)
		c.Header("X-Request-Id", id)

		c.Next()

		slog.Info("api.perf",
			slog.String("id", id),
			slog.Duration("dt", time.Since(t)),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.String("ip", c.ClientIP()),
		)
	}
}
