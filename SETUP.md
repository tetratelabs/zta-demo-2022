# NIST demo

Instructions are a WIP

## Requirements

* istioctl
* envsubst
* openssl
* kustomize

## Deployment steps

```bash
# Export env variables
export HUB=<CHANGEME>
export WASM_HUB=<CHANGEME>
export OIDC_CLIENT_SECRET=<CHANGEME>
export OIDC_CLIENT_ID=<CHANGEME>

# Install Istio
istioctl install -f istio/istio-operator.yaml

# Deploy all applications
kustomize build k8s/ | envsubst | kubectl apply -f -

# Generate the certificates for the Ingress gateway and
# expose the vulnerable app in the Istio ingress
bash istio/gen-certs.sh
kubectl -n nist-demo-2022 apply -f istio/vulnerable-app.yaml

# Enforce OIDC at the ingress level
kubectl apply -f istio/oidc-policy.yaml

# Install the WASM Plugin
envsubst < istio/wasm-patch.yaml | kubectl apply -f -
```
