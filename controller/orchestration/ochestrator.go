package orchestration

import (
	"context"

	"github.com/morty-faas/morty/controller/types"
)

type Orchestrator interface {
	// GetFunctions retrieve all the functions currently provisioned into the orchestrator.
	GetFunctions(ctx context.Context) ([]*types.Function, error)

	// CreateFunction register the function into the orchestrator, but doesn't deploy an instance of it.
	CreateFunction(ctx context.Context, fn *types.Function) (*types.Function, error)

	// GetFunctionInstance retrieve an instance of the function, that must be ready to receive requests.
	GetFunctionInstance(ctx context.Context, fn *types.Function) (*types.FnInstance, bool, error)

	// DeleteFunctionInstance delete a function instance.
	DeleteFunctionInstance(ctx context.Context, fn *types.FnInstance) error
}
