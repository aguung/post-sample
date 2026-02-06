package middleware

import (
	"fmt"
	"net/http"
	"time"

	"post/internal/pkg/logger"
	"post/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New()
		c.Set("RequestID", uuid.String())
		c.Header("X-Request-ID", uuid.String())
		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		reqID, _ := c.Get("RequestID")

		log := logger.GetLogger()
		log.Info().
			Str("request_id", reqID.(string)).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Str("ip", c.ClientIP()).
			Dur("duration", duration).
			Msg("Request processed")
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				reqID, _ := c.Get("RequestID")
				log.Error().
					Str("request_id", reqID.(string)).
					Interface("error", err).
					Msg("Panic recovered")

				response.Error(c, http.StatusInternalServerError, "Internal Server Error", fmt.Sprintf("%v", err))
				c.Abort()
			}
		}()
		c.Next()
	}
}
