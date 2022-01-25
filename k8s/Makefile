
install/%:
	istioctl install -y -f manifests/istio-operator.yaml
	cd $(@F) && bash setup.sh
	kustomize build manifests/ | envsubst | kubectl apply -f -
	envsubst < manifests/vulnerable-app-ingress.yaml | kubectl apply -f -

uninstall/%: unconfigure
	kustomize build manifests/ | envsubst | kubectl delete --ignore-not-found -f -
	cd $(@F) && bash uninstall.sh
	istioctl x uninstall --purge -y
	kubectl delete namespace istio-system --ignore-not-found

unconfigure:
	kubectl delete --ignore-not-found -f ../config/oidc-policy.yaml
	kubectl delete --ignore-not-found -f ../config/runtime-authn.yaml
	envsubst < ../config/wasm-patch.yaml | kubectl delete --ignore-not-found -f -
