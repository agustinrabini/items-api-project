package handlers

import (
	"github.com/agustinrabini/items-api-project/src/main/api/config"

	"github.com/gin-gonic/gin"
	"github.com/jopitnow/go-jopit-toolkit/goutils/logger"
)

func LoggerHandler(requestName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqLogger := logger.NewRequestLogger(c, requestName, config.LogRatio, config.LogBodyRatio)
		c.Next()
		reqLogger.LogResponse(c)
	}
}
