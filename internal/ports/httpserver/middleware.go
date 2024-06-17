package httpserver

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "log/slog"
    "time"
)

func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := uuid.New()
        c.Set("requestID", requestID)
        c.Next()
    }
}

func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        logger.Info(
            "got HTTP request",
            slog.String("requestID", c.GetString("requestID")),
            slog.String("method", c.Request.Method),
            slog.String("path", c.Request.URL.Path),
            slog.String("queryParams", c.Request.URL.RawQuery),
            slog.String("clientIP", c.ClientIP()),
        )

        c.Next()

        logger.Info(
            "request is handled",
            slog.String("requestID", c.GetString("requestID")),
            slog.Duration("latency", time.Since(start)),
            slog.Int("status", c.Writer.Status()),
            slog.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
        )
    }
}
