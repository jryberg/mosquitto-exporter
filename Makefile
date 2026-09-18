PKG_NAME:=github.com/jryberg/mosquitto-exporter
BUILD_DIR:=bin
MOSQUITTO_EXPORTER_BINARY:=$(BUILD_DIR)/mosquitto_exporter
IMAGE := jryberg/mosquitto-exporter
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS=-s -w -X main.Version=$(VERSION)
.PHONY: help
help:
	@echo
	@echo "Available targets:"
	@echo "  * build             - build the binary, output to $(MOSQUITTO_EXPORTER_BINARY)"
	@echo "  * linux             - build the binary, output to $(MOSQUITTO_EXPORTER_BINARY)"
	@echo "  * docker            - build docker image"

.PHONY: build
build: export CGO_ENABLED=0
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(MOSQUITTO_EXPORTER_BINARY) -ldflags="$(LDFLAGS)" $(PKG_NAME)

linux: export GOOS=linux
linux: build

docker:
	docker build -t $(IMAGE):$(VERSION) .
	docker build -t $(IMAGE):latest .

push:
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest
