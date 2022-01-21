# ZTA and DevSecOps for Cloud Native Applications demo

Instructions to deploy the ZTA demo in a Kubernetes cluster.

## Requirements

* istioctl 1.12
* envsubst
* openssl
* kustomize

## Deployment steps

make install/gke
make install/local

make uninstall/gke
make uninstall/local

kubectl run tmp-shell --rm -i --tty --image nicolaka/netshoot -- /bin/bash