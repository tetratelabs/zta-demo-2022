
install:
	istioctl install -y -f istio-operator.yaml
	kustomize build . | kubectl apply -f -

uninstall: unconfigure
	@# Restore the original Istio config, without NGAC
	istioctl install -y -f ../k8s/manifests/istio-operator.yaml
	kustomize build . | kubectl delete --ignore-not-found -f -

unconfigure:
	kubectl delete --ignore-not-found -f ../config/ngac-policy.yaml
