package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morty-faas/morty/controller/state"
)

type healthcheckResponse struct {
	Status string `json:"status"`
}

func HealthHandler(state state.State) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Currently, we consider that the server is always healthy
		// TODO: check for state health (and underlying orchestrator ?)
		c.JSON(http.StatusOK, &healthcheckResponse{
			Status: "UP",
		})
	}
}
