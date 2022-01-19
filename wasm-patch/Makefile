NAME := log4shell-patch
WASM := $(NAME).wasm

HUB := gcr.io/ignasi-nist-2022
TAG ?= 0.1

# Envoy >= 1.21.0 is required
ENVOY_IMAGE := envoyproxy/envoy:v1.21.0

OUT_DIR := build
OUT := $(WASM)

.PHONY: compile
compile: $(OUT)

$(OUT):
	tinygo build -o $(OUT) -scheduler=none -target=wasi ./...

.PHONY: test
test:
	go test -v -tags=proxytest ./...

.PHONY: clean
clean:
	rm -f $(OUT)

.PHONY: docker-build
docker-build: $(OUT)
	docker build --build-arg WASM_BINARY_PATH=$(OUT) -t $(HUB)/$(NAME):$(TAG) .

.PHONY: example-run
example-run: $(OUT)
	docker run --rm -p 8080:8080 \
		-v $$(pwd)/envoy-example.yaml:/opt/log4shell/envoy.yaml \
		-v $$(pwd)/$(OUT):/opt/log4shell/$(WASM) \
		$(ENVOY_IMAGE) -c /opt/log4shell/envoy.yaml
