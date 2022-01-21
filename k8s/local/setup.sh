source ../variables.env

set -e

DAYS=1825   # 5 years
CN=localhost
OUTDIR=$(mktemp -d)

pushd .
mkdir -p ${OUTDIR}
cd ${OUTDIR}

echo "Generating the CA..."
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -sha256  -out ca.crt -days ${DAYS} -subj "/C=US/ST=California/L=San Francisco/O=Tetrate"

echo "Generating the server certificate..."
openssl genrsa -out cert.key 2048
openssl req -new -key cert.key -out cert.csr -subj "/C=US/ST=California/L=San Francisco/O=Tetrate/CN=${CN}"
openssl x509 -req -extfile <(printf "subjectAltName=DNS:${CN}") -in cert.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out cert.crt -days ${DAYS} -sha256

kubectl delete secret -n istio-system ingress-tls-cert --ignore-not-found
kubectl create secret tls -n istio-system ingress-tls-cert --key=cert.key --cert=cert.crt

popd
