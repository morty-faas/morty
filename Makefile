BUILD_DIR:=.

default: build

.PHONY: build
build: build/controller build/registry

.PHONY: build/controller
build/controller:
	go build -o $(BUILD_DIR)/morty-controller cmd/controller/main.go

.PHONY: build/registry
build/registry:
	go build -o $(BUILD_DIR)/morty-registry cmd/registry/main.go
