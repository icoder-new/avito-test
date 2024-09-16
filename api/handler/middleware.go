package handler

import (
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

func NewMWLogger(log logger.ILogger) gin.HandlerFunc {
	log = logger.With(log, logger.String("component", "middleware/logger"))

	log.Info("logger middleware enabled")

	return func(c *gin.Context) {
		entry := logger.With(
			log,
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.String("remote_addr", c.ClientIP()),
			logger.String("user_agent", c.Request.UserAgent()),
			logger.String("request_id", c.GetHeader("X-Request-Id")),
		)

		startTime := time.Now()

		c.Next()

		entry.Info("request completed",
			logger.Int("status", c.Writer.Status()),
			logger.Int("bytes", c.Writer.Size()),
			logger.String("duration", time.Since(startTime).String()),
		)
	}
}
