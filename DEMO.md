# ZTA and DevSecOps for Cloud Native Applications demo

The following steps show how a Service Mesh can leverage the
features of a Security Kernel suitable for a Zero Trust Architecture
Platform. We will see:

* How a service mesh leverages **runtime identities* to protect
  service-to-service communications.
* How it can enforce **user identity based policies** and integrate with
  external or corporate Identity Providers.
* How policy is **enforced at the application level**, not only at the
  perimeter.
* How **targeted application policies** can be applied to affect only a subset
  of the applications in the mesh.

The demo consists of a Java application that is vulnerable to the `Log4Shell`
exploit. We will use the service mesh to enforce that access is authenticated on the
corporate Identity Provider, that only the right users and local services can access
the application, and that malicious payloads that trigger the `Log4Shell` exploit
are rejected.

## Install the demo environment

Before starting, install the demo environment following the instructions
in the [k8s/SETUP.md](k8s/SETUP.md) file.

## 1. Access the vulnerable application

Once you have deployed the demo environment you will be able to access the vulnerable
application at:

* `https://<your app DNS name>` - if you used the `gke` mode.
* `https://localhost:8443` - if you used the `local` mode.

You will see something like this:

```
Welcome, anonymous!
```

## 2. Enforce user identities

We don't want to allow unauthenticated users to access our application, so let's
apply a policy that configures the ingress gateway to require authentication against
the corporate identity Provider:

```
$ kubectl apply -f config/oidc-policy.yaml
```

If you inspect the policy file you'll see that it applies to the `istio-ingressgateway`
and that it uses a `CUSTOM` target that delegates to the `authservice-grpc` provider.
The provider is configured in The Istio global mesh config that you can check with:

```
$ kubectl -n istio-system describe configmap istio
```

Once the policy is applied you can refresh the browser and you will be redirected to the
Identity Provider login page. After accepting in the consent screen and logging in, you'll
see something like:

```
Welcome, Ignasi!


Authenticated with token:
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkRGRWVyODZqY2lRQTNfUVdETkE3MyJ9.eyJuaWNrbmFtZSI6ImlnbmFzaSt0ZXN0IiwibmFtZSI6IklnbmFzaSIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8wNDRjMjU1MDk4NDE2M2M5OTQ3N2NkM2QyYjY0NWViND9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmlnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIyLTAxLTIzVDEwOjQyOjA2Ljg2NVoiLCJlbWFpbCI6ImlnbmFzaSt0ZXN0QHRldHJhdGUuaW8iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiaXNzIjoiaHR0cHM6Ly9uYWN4LWRtei5ldS5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NjFlN2Y0NWE3NmRjM2EwMDZhYWU1MGRiIiwiYXVkIjoiZHl5VzBtZTRMcThuM3RZMzBGYXR1RFFGWHB0WnZtNG8iLCJpYXQiOjE2NDI5MzQ1MjcsImV4cCI6MTY0Mjk3MDUyNywibm9uY2UiOiJsOUV4dW1EUUZQc3Z4VzJ1YkhSUHpiUDA1aFJuYnQ1dFBKNXJaRXVrSFlvIn0.AUr-S54HmRosQaHLFN99hxj1eP1NyDkk_42Bihlbh5OyQdTba-J_KYwWqUHPYdry8RQmHzYz6moDH3hynV5TRCDP0TNzCc9Y6eYoQmT0U59ZeKL1d38XMmGdTBTbzwMCRKoGf4wyopBPsFsIE2tH4iUMBiL7uKw_0kgjOr2UxJcUQR7bPVyvRSXIanxdrrtSWgpBdibAZ80c-z2V7m9uWJM8Tz_SjVwVm1PiXh_nkptFnWfq-i8J_aMNpDLU_RIn2D4nz_omBLdvSRQYjKyFMUIQZ8Huctpx-bKcBJFHT7l2QFjFMyaB00GKFLbsFCk3EDNeWQ8E8a8nn4fPooabbg
```

## 3. Enforce runtime access

Now we have configured our ingress to require an user login. We could have applied this policy
directly to the application sidecar, but for the demo purposes we'll do it just at the ingress level.

This means that any other service in the cluster could directly reach the application withot going through
the ingress. We can check it by launching a new pod and accessing the app as follows:

```
$ kubectl run tmp-shell --rm -i --tty --image nicolaka/netshoot -- /bin/bash
bash-5.1# curl http://vulnerable:8080
Welcome, anonymous!
bash-5.1# exit
```

Let's create a runtime policy that enforces that our application can only be reached from the ingress gateway.
Services in the cluster will no longer have direct access to it:

```
$ kubectl apply -f config/runtime-authn.yaml
```

If you inspect the contents of the policy you'll see that it applies to the `vulnerable` application and that
it only allows access from a specific source principal. That source principal matches the
[SPIFEE identity](https://spiffe.io/docs/latest/spiffe-about/spiffe-concepts/#spiffe-id) of the Istio Ingress Gateway.

We can now try to directly access the application again from inside the cluster:

```
$ kubectl run tmp-shell --rm -i --tty --image nicolaka/netshoot -- /bin/bash
bash-5.1# curl http://vulnerable:8080
RBAC: access denied
bash-5.1# exit
```

Now we get an access denied, because the proxy sidecar in the application pod is rejecting the connection since the
runtime identity presented by our workload does not match the configured one.
