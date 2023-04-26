BUILD_DIR:=bin

default: build

.PHONY: build
build: build/controller build/registry

.PHONY: build/controller
build/controller:
	go build -o $(BUILD_DIR)/controller cmd/controller/main.go

.PHONY: build/registry
build/registry:
	go build -o $(BUILD_DIR)/registry cmd/registry/main.go
