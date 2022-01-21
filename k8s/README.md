# ZTA and DevSecOps for Cloud Native Applications demo

The demo scenario is deployed on a GKE Kubernetes cluster and uses
[Auth0](https://auth0.com/) as an OIDC Identity Provider. Follow these
instructions to deploy the demo scenario on a Kubernetes cluster.

## Requirements

The following tools are required to build the environment:

* [istioctl 1.12](https://istio.io/latest/docs/setup/getting-started/#download)
* gcloud
* envsubst
* openssl
* kustomize

## Demo profiles

The demo can be deployed in two different ways depending on the infrastructure that is available:

| Profile | Descrption |
|---------|------------|
| local   | Deploys everything on the GKE cluster and relies on Kubernetes port-forwarding to expose the Istio Ingress Gateway in `localhost`. |
| gke | Uses `cert-manager` and `Google Cloud DNS` to configure the certificates and hostnames where the application will be exposed via the Istio Ingress Gateway. Requires a real DNS zone managed in the GKE project. |

## Auth0 application setup

The demo uses [Auth0](https://auth0.com/) as an OIDC Identity Provider to enforce users are authenticated.
Follow these instructions to configure an application for the demo:

* In the Auth0 admin console, go to **Applications > Applications > Create Application**.
* Select **Regular Web Applications**, give it a name, and **Create**.
* Go to the **Settings** tab and configure the following, according to the profile you want to use (gke/local):
  * **Allowed Callback URLs:**
    * For the `local` profile: `https://localhost:8443/oauth/callback`
    * For the `gke` profile: `https://<desired app DNS name>/oauth/callback`
  * **Allowed Logout URLs:**
    * For the `local` profile: `https://localhost:8443/logout`
    * For the `gke` profile: `https://<desired app DNS name>/logout`
  * Scroll down to the **Advanced** section, and in the **OAuth** tag, make sure the
    **OIDC Conformant** option is enabled.
* Go to **Users Management > Users > Create User**. Enter the requested data and **Create**.

## Environment setup

Once the Auth0 application has been created, configure the environment variables that
will be used to customize the deployment process. Make a copy of the
[variables.example.env](variables.example.env) file, name it `variables.env`, and modify
it according to your needs.

Once the file is there, you're ready to go!

## Deployment steps

Before applying the Kubernetes manifest following the instructions below, build all the images:

```
source variables.env
make -C ../ clean docker-build docker-push
```

### gke mode

To deploy the demo in `gke` mode and leverage automatic certificate issuance and DNS
configuration, do it as follows:

```
make install/gke
```

Once the deployment completes you can open a browser to `https://<your app DNS name>`

### local mode

To deploy the demo in `local` mode, install the application as follows and expose the
Istio ingress gateway locally:

```
make install/local
kubectl -n istio-system port-forward svc/istio-ingressgateway 8443:443
```

Once the deployment completes you can open a browser to `https://localhost:8443`

## Environment cleanup

To cleanup the environment and uninstall everything from the cluster, you can use
the following commands, according to the mode you used in the installation process:

```
make uninstall/gke
```
or
```
make uninstall/local
```
