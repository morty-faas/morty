package config

import (
	"github.com/morty-faas/morty/controller/orchestration"
	"github.com/morty-faas/morty/controller/orchestration/rik"
	"github.com/morty-faas/morty/controller/state"
	"github.com/morty-faas/morty/controller/state/memory"
	"github.com/morty-faas/morty/controller/state/redis"
	log "github.com/sirupsen/logrus"
	"github.com/thomasgouveia/go-config"
)

type (
	Config struct {
		Port         int          `yaml:"port"`
		MetricsPort  int          `yaml:"metricsPort"`
		Orchestrator Orchestrator `yaml:"orchestrator"`
		State        State        `yaml:"state"`
	}

	Orchestrator struct {
		Rik rik.Config `yaml:"rik"`
	}

	State struct {
		Redis redis.Config `yaml:"redis"`
	}
)

var loaderOptions = &config.Options[Config]{
	Format: config.YAML,

	// Environment variables lookup
	EnvEnabled: true,
	EnvPrefix:  "MORTY_CONTROLLER",

	// Configuration file
	FileName:      "controller",
	FileLocations: []string{"/etc/morty", "$HOME/.morty", "."},

	// Default configuration
	Default: &Config{
		Port:        8080,
		MetricsPort: 9090,
		Orchestrator: Orchestrator{
			Rik: rik.Config{
				Cluster: "http://localhost:5000",
			},
		},
	},
}

// Load the configuration from the different sources (environment, files, default)
func Load() (*Config, error) {
	cl, err := config.NewLoader(loaderOptions)
	if err != nil {
		return nil, err
	}

	cfg, err := cl.Load()
	if err != nil {
		return nil, err
	}

	log.Debugf("Loaded configuration: %+v", cfg)

	return cfg, nil
}

// StateFactory initializes a new state implementation based on the configuration.
func (c *Config) StateFactory(expiryCallback state.FnExpiryCallback) (state.State, error) {
	log.Debugf("Applying state factory based on configuration")
	if err := ensureKeyHasSingleSubKey(c.State); err != nil {
		return nil, err
	}

	if isDefined(c.State.Redis) {
		return redis.NewState(&c.State.Redis, expiryCallback)
	}

	// By default, we will use a in memory state engine if no configuration
	// is provided by the user.
	return memory.NewState(), nil
}

// OrchestratorFactory initializes a new orchestrator implementation based on the configuration.
func (c *Config) OrchestratorFactory() (orchestration.Orchestrator, error) {
	log.Debugf("Applying orchestrator factory based on configuration")
	return rik.NewOrchestrator(&c.Orchestrator.Rik)
}
