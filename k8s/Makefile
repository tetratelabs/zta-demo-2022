
install/%:
	istioctl install -y -f manifests/istio-operator.yaml
	cd $(@F) && bash setup.sh
	kustomize build manifests/ | envsubst | kubectl apply -f -

uninstall/%:
	kubectl delete --ignore-not-found -f ../config/oidc-policy.yaml
	envsubst < ../config/wasm-patch.yaml | kubectl delete --ignore-not-found -f -
	kustomize build manifests/ | envsubst | kubectl delete --ignore-not-found -f -
	cd $(@F) && bash uninstall.sh
	istioctl x uninstall --purge -y
	kubectl delete namespace istio-system --ignore-not-found
