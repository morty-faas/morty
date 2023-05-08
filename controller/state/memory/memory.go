package memory

import (
	"context"
	"time"

	"github.com/morty-faas/morty/controller/state"
	"github.com/morty-faas/morty/controller/types"
	log "github.com/sirupsen/logrus"
)

// adapter is an implementation of the state.State interface
type adapter struct {
	versions map[string][]string
	store    map[string]*types.Function
}

var _ state.State = (*adapter)(nil)

// NewState initializes a new state adapter for Memory engine.
func NewState() state.State {
	log.Info("State engine 'memory' successfully initialized")
	return &adapter{
		versions: make(map[string][]string),
		store:    make(map[string]*types.Function),
	}
}

func (a *adapter) GetAll(ctx context.Context) ([]*types.Function, error) {
	var functions []*types.Function
	for _, fn := range a.store {
		functions = append(functions, fn)
	}
	return functions, nil
}

func (a *adapter) Get(ctx context.Context, key string) (*types.Function, error) {
	return a.GetByVersion(ctx, key, state.DefaultVersion)
}

func (a *adapter) GetByVersion(ctx context.Context, key string, version string) (*types.Function, error) {
	log.Tracef("state/memory: retrieving value for key '%s' and version '%s'", key, version)
	versions := a.getFunctionVersions(key)

	if version == state.DefaultVersion {
		version = versions[len(versions)-1]
	}

	v, exists := a.store[key+":"+version]
	if !exists {
		return nil, state.ErrKeyNotFound
	}

	return v, nil
}

func (a *adapter) Set(ctx context.Context, fn *types.Function) error {
	log.Tracef("state/memory: setting value '%+v' for key '%s'", fn, fn.Name)
	a.versions[fn.Name] = append(a.getFunctionVersions(fn.Name), fn.Version)
	a.store[fn.Id()] = fn
	return nil
}

func (a *adapter) SetMultiple(ctx context.Context, functions []*types.Function) []error {
	var errors []error
	for _, fn := range functions {
		if err := a.Set(ctx, fn); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (a *adapter) SetWithExpiry(ctx context.Context, key string, expiry time.Duration) error {
	// Expiration isn't supported for memory engine
	return nil
}

func (a *adapter) getFunctionVersions(key string) []string {
	versions, exists := a.versions[key]
	if !exists {
		return []string{}
	}
	return versions
}
