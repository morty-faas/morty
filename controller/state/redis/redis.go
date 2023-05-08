package redis

import (
	"context"
	"strings"
	"time"

	"github.com/morty-faas/morty/controller/state"
	"github.com/morty-faas/morty/controller/types"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

// adapter is an implementation of the state.State interface
type adapter struct {
	client *redis.Client
}

// Config hold the configuration about the Redis state adapter
type Config struct {
	Addr string `yaml:"addr"`
}

var _ state.State = (*adapter)(nil)

// NewState initializes a new state adapter for Redis based on the given configuration.
// An error could be returned if any errors happens during the adapter initialization.
func NewState(cfg *Config, expiryCallback state.FnExpiryCallback) (state.State, error) {
	log.Debugf("Bootstrapping Redis state adapter with options: %#v", cfg)
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		DB:   0,
	})

	// Enable Keyspace events as we will need them to handle function instances expiration
	if _, err := client.ConfigSet(context.Background(), "notify-keyspace-events", "KEA").Result(); err != nil {
		log.Errorf("Failed to enable Redis Keyspace Events: %v", err)
		return nil, err
	}

	// this is telling redis to subscribe to events published in the keyevent channel, specifically for expired events
	pubsub := client.PSubscribe(context.Background(), "__keyevent@0__:expired")

	go func(*redis.PubSub) {
		for {
			message, err := pubsub.ReceiveMessage(context.Background())
			if err != nil {
				log.Errorf("failed to receive expiry event: %v", err)
				continue
			}
			log.Tracef("Key %s has expired", message.Payload)
			expiryCallback(message.Payload)
		}
	}(pubsub)

	log.Info("State engine 'redis' successfully initialized")
	return &adapter{client}, nil
}

func (a *adapter) GetAll(ctx context.Context) ([]*types.Function, error) {
	var functions []*types.Function

	keys, err := a.client.Keys(ctx, "*:v*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		fn, err := a.getFunctionByKey(ctx, key)
		if err != nil {
			log.Warnf("state/redis: failed to get function at key '%s' : %v", key, err)
		}
		functions = append(functions, fn)
	}

	return functions, nil
}

func (a *adapter) GetByVersion(ctx context.Context, key, version string) (*types.Function, error) {
	log.Tracef("state/redis: retrieving value for key '%s' and version '%s'", key, version)
	versions, err := a.getFunctionVersions(ctx, key)
	if err != nil {
		return nil, err
	}

	if version == state.DefaultVersion {
		version = versions[len(versions)-1]
	}

	return a.getFunctionByKey(ctx, key+":"+version)
}

func (a *adapter) Get(ctx context.Context, key string) (*types.Function, error) {
	return a.GetByVersion(ctx, key, state.DefaultVersion)
}

func (a *adapter) Set(ctx context.Context, fn *types.Function) error {
	r := a.client.HSet(ctx, fn.Id(), fn)
	log.Tracef("state/redis: %s", r.String())
	_, err := r.Result()
	if err != nil {
		return err
	}

	// Maintain the list of versions per function
	r = a.client.RPush(ctx, fn.Name, fn.Version)
	log.Tracef("state/redis: %s", r.String())
	_, err = r.Result()
	return err
}

func (a *adapter) SetMultiple(ctx context.Context, functions []*types.Function) []error {
	errors := []error{}
	for _, fn := range functions {
		if err := a.Set(ctx, fn); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (a *adapter) SetWithExpiry(ctx context.Context, key string, expiry time.Duration) error {
	log.Debugf("Set expiration of %v for key %s", expiry, key)

	_, err := a.client.Set(ctx, key, "", expiry).Result()
	return err
}

func (a *adapter) getFunctionByKey(ctx context.Context, key string) (*types.Function, error) {
	r := a.client.HGetAll(ctx, key)
	log.Tracef("state/redis: %s", r.String())
	res, err := r.Result()
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	tokens := strings.Split(key, ":")
	name, version := tokens[0], tokens[1]

	fn := &types.Function{
		Name:           name,
		Version:        version,
		OrchestratorId: res["orchestratorId"],
		ImageURL:       res["imageUrl"],
	}

	return fn, nil
}

func (a *adapter) getFunctionVersions(ctx context.Context, key string) ([]string, error) {
	return a.client.LRange(ctx, key, 0, -1).Result()
}
