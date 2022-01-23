source ../variables.env

# Install cert-manager
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml

# Make sure the Cloud DNS API is enabled
gcloud services enable dns.googleapis.com containerregistry.googleapis.com --project=${PROJECT_ID}

# Configure the service account and secret cert-manager will use
# to solve the ACME challenges presented by Let's Encrypt when issuing
# certificates.
# This account will be used by external-dns as well to automatically create
# DNS records for the hostnames exposed in the mesh.
gcloud iam service-accounts create dns01-solver \
    --project ${PROJECT_ID} \
    --display-name "dns01-solver"

gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member serviceAccount:dns01-solver@${PROJECT_ID}.iam.gserviceaccount.com \
    --role roles/dns.admin

gcloud iam service-accounts keys create /tmp/key.json \
    --iam-account dns01-solver@${PROJECT_ID}.iam.gserviceaccount.com

kubectl delete secret clouddns-dns01-solver-sa --ignore-not-found
kubectl create secret generic clouddns-dns01-solver-sa \
    -n cert-manager \
    --from-file=/tmp/key.json

envsubst < external-dns.yaml | kubectl apply -n cert-manager -f -

echo "Waiting a bit for cert-manager..."
sleep 60

envsubst < issuer.yaml | kubectl apply -f -
envsubst < certificate.yaml | kubectl apply -f -
