-include ../common.mk

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
	docker build --build-arg WASM_BINARY_PATH=$(OUT) -t $(HUB)/$(NAME):$(TAG) .

docker-push:
	docker push $(HUB)/$(NAME):$(TAG)
