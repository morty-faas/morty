package config

import (
	"github.com/thomasgouveia/go-config"
)

type (
	Config struct {
		Port    int     `yaml:"port"`
		Storage Storage `yaml:"storage"`
	}

	Storage struct {
		S3 S3 `yaml:"s3"`
	}

	S3 struct {
		Bucket   string `yaml:"bucket"`
		Region   string `yaml:"region"`
		Endpoint string `yaml:"endpoint"`
	}
)

var loaderOptions = &config.Options[Config]{
	Format: config.YAML,

	// Configure the loader to lookup for environment
	// variables with the following pattern: APP_*
	EnvEnabled: true,
	EnvPrefix:  "app",

	// Configure the loader to search for an alpha.yaml file
	// inside one or more locations defined in `FileLocations`
	FileName:      "registry",
	FileLocations: []string{"/etc/morty-registry", "$HOME/.morty", "."},

	// Inject a default configuration in the loader
	Default: &Config{
		Port: 8081,
		Storage: Storage{
			S3: S3{
				Bucket:   "functions",
				Region:   "eu-west-1",
				Endpoint: "http://localhost:9000",
			},
		},
	},
}

// Load the configuration from the different sources (default, file and environment).
func Load() (*Config, error) {
	cl, err := config.NewLoader(loaderOptions)
	if err != nil {
		return nil, err
	}
	return cl.Load()
}
