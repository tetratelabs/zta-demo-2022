kubectl delete secret ingress-tls-cert -n istio-system --ignore-not-found
kubectl delete secret clouddns-dns01-solver-sa -n cert-manager --ignore-not-found

envsubst < external-dns.yaml | kubectl delete --ignore-not-found -f -
envsubst < issuer.yaml | kubectl delete --ignore-not-found -f -
envsubst < certificate.yaml | kubectl delete --ignore-not-found -f -

kubectl delete -f https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml

kubectl delete namespace cert-manager --ignore-not-found

#gcloud iam service-accounts delete dns01-solver@${PROJECT_ID}.iam.gserviceaccount.com --quiet
