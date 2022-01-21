# Common variables

# Change to push application images to a different registry
export HUB ?= gcr.io/ignasi-nist-2022

# Istio 1.12 does not yet support imagePullSecrets in
# WasmPlugins to pull images from private registries.
# We need to push the WASM patch to a public Docker registry.
export WASM_HUB ?= docker.io/nacx

# Tag used for all the iamges
export TAG ?= 0.1
