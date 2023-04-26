package memory

import (
	"context"
	"errors"
	"time"

	"github.com/morty-faas/morty/controller/state"
	"github.com/morty-faas/morty/controller/types"
	log "github.com/sirupsen/logrus"
)

// adapter is an implementation of the state.State interface
type adapter struct {
	store map[string]*types.Function
}

var _ state.State = (*adapter)(nil)

// NewState initializes a new state adapter for Memory engine.
func NewState() state.State {
	log.Info("State engine 'memory' successfully initialized")
	return &adapter{
		store: make(map[string]*types.Function),
	}
}

func (a *adapter) Get(ctx context.Context, key string) (*types.Function, error) {
	log.Tracef("state/memory: retrieving value for key '%s'", key)
	v, exists := a.store[key]
	if !exists {
		return nil, state.ErrKeyNotFound
	}
	return v, nil
}

func (a *adapter) Set(ctx context.Context, fn *types.Function) error {
	log.Tracef("state/memory: setting value '%+v' for key '%s'", fn.Id, fn)
	a.store[fn.Name] = fn
	return nil
}

func (a *adapter) SetMultiple(ctx context.Context, functions []*types.Function) []error {
	for _, fn := range functions {
		a.Set(ctx, fn)
	}
	return nil
}

func (a *adapter) SetWithExpiry(ctx context.Context, key string, expiry time.Duration) error {
	return errors.New("not supported")
}
