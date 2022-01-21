# Common variables

export HUB ?= gcr.io/ignasi-nist-2022

# Istio's WasmPlugin does not yet support imagePullSecrets
# to pull images from private registries
export WASM_HUB ?= docker.io/nacx

export TAG ?= 0.1
