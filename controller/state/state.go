package state

import (
	"context"
	"errors"
	"time"

	"github.com/morty-faas/morty/controller/types"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

const DefaultVersion = "latest"

type FnExpiryCallback func(string)

// State is a generic interface for our functions state
type State interface {
	// GetAll retrieve all the functions from the state.
	GetAll(ctx context.Context) ([]*types.Function, error)
	// Get retrieve the latest version of the function associated to the given key.
	// If the key doesn't exists, an error ErrKeyNotFound will be returned
	Get(ctx context.Context, key string) (*types.Function, error)
	// GetByVersion retrieve a function associated to the given version.
	GetByVersion(ctx context.Context, key, version string) (*types.Function, error)
	// Set a tuple of key/value in the state
	Set(ctx context.Context, fn *types.Function) error
	// SetMultiple set multiple keys in one call
	SetMultiple(ctx context.Context, functions []*types.Function) []error

	SetWithExpiry(ctx context.Context, key string, expiry time.Duration) error
}
