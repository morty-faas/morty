package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morty-faas/morty/controller/orchestration"
	"github.com/morty-faas/morty/controller/state"
	"github.com/morty-faas/morty/controller/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"
)

type createFnRequest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Image   string `json:"image"`
}

var (
	ErrVersionConflict = errors.New("version already exists for this function")
)

func CreateFunctionHandler(state state.State, orch orchestration.Orchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Parse the request body
		data := &createFnRequest{}
		if err := c.BindJSON(data); err != nil {
			logrus.Errorf("Failed to decode create function request body: %v", err)
			c.JSON(http.StatusBadRequest, makeApiError(err))
			return
		}

		if !semver.IsValid(data.Version) {
			err := fmt.Errorf("Version '%s' isn't a valid semantic version", data.Version)
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, makeApiError(err))
		}

		fn := &types.Function{
			Name:     data.Name,
			Version:  data.Version,
			ImageURL: data.Image,
		}

		// Ensure that the version do not exists for this function
		if fn, _ := state.GetByVersion(ctx, fn.Name, fn.Version); fn != nil {
			c.JSON(http.StatusConflict, makeApiError(ErrVersionConflict))
			return
		}

		// Process the function into the orchestrator
		fn, err := orch.CreateFunction(ctx, fn)
		if err != nil {
			logrus.Errorf("Failed to create function into the orchestrator: %v", err)
			c.JSON(http.StatusInternalServerError, makeApiError(err))
			return
		}

		if err := state.Set(ctx, fn); err != nil {
			logrus.Errorf("Failed to add function to the state: %v", err)
			c.JSON(http.StatusInternalServerError, makeApiError(err))
			return
		}

		c.JSON(http.StatusOK, fn)
	}
}
