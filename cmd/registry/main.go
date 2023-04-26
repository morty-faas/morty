package main

import (
	"os"
	"os/user"

	"github.com/morty-faas/morty/registry/registry"
	log "github.com/sirupsen/logrus"
)

const (
	logEnvVarKey = "REGISTRY_LOG"
)

func main() {
	// Init logger for the app
	level := log.InfoLevel
	envLevel := os.Getenv(logEnvVarKey)
	if envLevel != "" {
		lvl, err := log.ParseLevel(envLevel)
		if err != nil {
			log.Fatalf("failed to parse log level from environment: %v", err)
		}
		level = lvl
	}
	log.SetLevel(level)

	// Check that the app is running as root
	// We need to run as root as we make system calls that requires elevated privileges
	// on the host (such as mount, mkfs.ext4, ...)
	u, err := user.Current()
	if err != nil {
		log.Fatalf("failed to retrieve user: %v", err)
	}

	// The user is root if the uid is 0
	if u.Uid != "0" {
		log.Fatal("registry needs root privileges to work properly")
	}

	// Run the registry HTTP server
	reg, err := registry.NewServer()
	if err != nil {
		log.Fatalf("failed to initialize the registry: %v", err)
	}

	reg.Serve()
}
