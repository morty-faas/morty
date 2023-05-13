BUILD_DIR:=.

GIT_IMPORT="github.com/morty-faas/morty/build"
GIT_COMMIT=$$(git rev-parse --short HEAD)
GIT_TAG=$$(git describe --abbrev=0 --tags)

LDFLAGS="-s -w -X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT) -X $(GIT_IMPORT).Version=$(GIT_TAG)"

default: build

.PHONY: build
build: build/controller build/registry

.PHONY: build/controller
build/controller:
	go build -ldflags $(LDFLAGS) -o $(BUILD_DIR)/morty-controller cmd/controller/main.go

.PHONY: build/registry
build/registry:
	go build -ldflags $(LDFLAGS) -o $(BUILD_DIR)/morty-registry cmd/registry/main.go
