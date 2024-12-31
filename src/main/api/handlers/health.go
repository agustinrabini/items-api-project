package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckerHandler struct provides the handler for a health check endpoint.
type HealthCheckerHandler struct{}

func NewHealthCheckerHandler() HealthCheckerHandler {
	return HealthCheckerHandler{}
}

// Ping is the handler of test app
// @Summary Ping
// @Description test if the router works correctly SAPEEE
// @Tags ping
// @Produce  json
// @Success 200
// @Router /ping [get]
func (h HealthCheckerHandler) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
