package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morty-faas/morty/controller/orchestration"
	"github.com/morty-faas/morty/controller/state"
)

func ListFunctionsHandler(s state.State, orch orchestration.Orchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		functions, err := orch.GetFunctions(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, makeApiError(err))
			return
		}

		c.JSON(http.StatusOK, functions)
	}
}
