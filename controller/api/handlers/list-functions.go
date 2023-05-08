package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morty-faas/morty/controller/orchestration"
	"github.com/morty-faas/morty/controller/state"
)

type (
	listFunctionResponse []functionLite
	functionLite         struct {
		Name     string   `json:"name"`
		Versions []string `json:"versions"`
	}
)

func ListFunctionsHandler(s state.State, orch orchestration.Orchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		functions, err := s.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, makeApiError(err))
			return
		}

		mapping := map[string][]string{}
		for _, fn := range functions {
			versions, exists := mapping[fn.Name]
			if !exists {
				versions = []string{}
			}

			mapping[fn.Name] = append(versions, fn.Version)
		}

		response := listFunctionResponse{}
		for k, v := range mapping {
			response = append(response, functionLite{
				Name:     k,
				Versions: v,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}
