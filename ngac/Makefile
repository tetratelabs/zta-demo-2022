
install:
	istioctl install -y -f istio-operator.yaml
	kustomize build . | kubectl apply -f -
	kubectl apply -f ../config/ngac-policy.yaml

uninstall:
	kubectl delete -f ../config/ngac-policy.yaml
	@# Restore the original Istio config, without NGAC
	istioctl install -y -f ../k8s/manifests/istio-operator.yaml
	kustomize build . | kubectl delete --ignore-not-found -f -
