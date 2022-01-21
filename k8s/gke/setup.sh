source ../variables.env

# Configure the service account and secret cert-manager will use
# to solve the ACME challenges presented by Let's Encrypt when issuing
# certificates.
gcloud iam service-accounts create dns01-solver \
    --project ${PROJECT_ID} \
    --display-name "dns01-solver"

gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member serviceAccount:dns01-solver@${PROJECT_ID}.iam.gserviceaccount.com \
    --role roles/dns.admin

gcloud iam service-accounts keys create /tmp/key.json \
    --iam-account dns01-solver@${PROJECT_ID}.iam.gserviceaccount.com

# Install cert-manager
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml

kubectl create secret generic clouddns-dns01-solver-sa \
    -n cert-manager \
    --from-file=/tmp/key.json

envsubst < external-dns.yaml | kubectl apply -n cert-manager -f -

echo "Waiting a bit for cert-manager..."
sleep 30

envsubst < issuer.yaml | kubectl apply -f -
envsubst < certificate.yaml | kubectl apply -f -
