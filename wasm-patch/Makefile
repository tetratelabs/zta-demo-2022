-include ../common.mk

NAME := log4shell-patch
WASM := $(NAME).wasm

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
