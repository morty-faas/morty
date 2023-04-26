package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morty-faas/morty/controller/orchestration"
	"github.com/morty-faas/morty/controller/state"
	"github.com/morty-faas/morty/controller/types"
	"github.com/sirupsen/logrus"
)

type createFnRequest struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

var (
	ErrNameConflict = errors.New("a function already exists with the given name")
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

		// We don't want to allow the user to create a function
		// if a function already exists with the same name at the creation
		if fn, _ := state.Get(ctx, data.Name); fn != nil {
			logrus.Errorf("A function already exists with the name: %s", data.Name)
			c.JSON(http.StatusConflict, makeApiError(ErrNameConflict))
			return
		}

		fn := &types.Function{
			Name:     data.Name,
			ImageURL: data.Image,
		}

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
