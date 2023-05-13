package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morty-faas/morty/build"
)

type wellknownResponse struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
}

func WellknownHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, &wellknownResponse{
			GitCommit: build.GitCommit,
			Version:   build.Version,
		})
	}
}
