-include ../common.mk

# Istio's WasmPlugin does not yet support imagePullSecrets
# to pull images from private registries
WASM_HUB := docker.io/nacx

NAME := log4shell-patch
WASM := $(NAME).wasm

# Envoy >= 1.21.0 is required
ENVOY_IMAGE := envoyproxy/envoy:v1.21.0

OUT_DIR := build
OUT := $(WASM)

compile: $(OUT)

$(OUT):
	tinygo build -o $(OUT) -scheduler=none -target=wasi ./...

test:
	go test -v -tags=proxytest ./...

clean:
	rm -f $(OUT)

docker-build: $(OUT)
	docker build --platform linux/amd64 --build-arg WASM_BINARY_PATH=$(OUT) -t $(WASM_HUB)/$(NAME):$(TAG) .

docker-push:
	docker push $(WASM_HUB)/$(NAME):$(TAG)
