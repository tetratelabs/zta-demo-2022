# Policy enforcement with NGAC

In this section we'll see how we can enforce access based on an NGAC graph. We will do
that based on some custom claims that are present in the JWT token of the authenticated user.

## Enabling custom claims in Auth0

Before we begin, we need to configure Auth0 to include additional custom claims
in the issued tokens. we can do that as follows in the Auth0 managament console:

* Go to the User, edit it, and add the following in the `user_metadata` field, then **Save**:
  ```json
  {
    "group": "Engineering"
  }
  ```
  You can add all information you want, but this example will use the `group` claim.
* Now we have to create an _Action_ that will inject all user metadata as custom claims in the issued tokens.
  Go to **Actions > Library**. Click the **Build Custom** button, select the **Login / Post Login**, give it
  a name, and **Create**. In the next screen, paste this code snippet and click **Deploy**:
  ```javascript
  exports.onExecutePostLogin = async (event, api) => {
    const namespace = 'https://zta-demo/';
    for (let key in event.user.user_metadata) {
      api.idToken.setCustomClaim(namespace + key, event.user.user_metadata[key]);
    }
  };
  ```
  Once that is done, go to **Actions > Flows**. Select **Login**, drag & drop the custom action you've just
  created between the _Login_ and the _Complete_ boxes, and click **Apply**.

## Deploy the NGAC enforcer

Once the user has been configured with additional claims, let's deploy the NGAC enforcer as follows:

```bash
$ source ../k8s/variables.env
$ make install
$ kubectl apply -f ../config/ngac-policy.yaml
```

Similar to the OIDC policy, if you inspect [the policy](../config/ngac-policy.yaml) file you'll see
that it applies to the `vulnerable` app and that it uses a `CUSTOM` target that delegates to the
`ngac-grpc` provider.

## Group based access control

Before starting, make sure you logout by going to the `/logout` path so you get a new token with the claims
we just created.

After deploying the NGAC enforcer, you'll see that you still have access to the application. If
you open the browser and go to the app, you'll see something like:

```
Welcome, Ignasi!
Group: Engineering
Accessing: /

Authenticated with token:
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkRGRWVyODZqY2lRQTNfUVdETkE3MyJ9.eyJodHRwczovL3p0YS1kZW1vL2NvdW50cnkiOiJFUyIsImh0dHBzOi8venRhLWRlbW8vZ3JvdXAiOiJFbmdpbmVlcmluZyIsIm5pY2tuYW1lIjoiaWduYXNpK3Rlc3QiLCJuYW1lIjoiSWduYXNpIiwicGljdHVyZSI6Imh0dHBzOi8vcy5ncmF2YXRhci5jb20vYXZhdGFyLzA0NGMyNTUwOTg0MTYzYzk5NDc3Y2QzZDJiNjQ1ZWI0P3M9NDgwJnI9cGcmZD1odHRwcyUzQSUyRiUyRmNkbi5hdXRoMC5jb20lMkZhdmF0YXJzJTJGaWcucG5nIiwidXBkYXRlZF9hdCI6IjIwMjItMDEtMjRUMjE6NDU6MDMuMTE0WiIsImVtYWlsIjoiaWduYXNpK3Rlc3RAdGV0cmF0ZS5pbyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJpc3MiOiJodHRwczovL25hY3gtZG16LmV1LmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw2MWU3ZjQ1YTc2ZGMzYTAwNmFhZTUwZGIiLCJhdWQiOiJkeXlXMG1lNExxOG4zdFkzMEZhdHVEUUZYcHRadm00byIsImlhdCI6MTY0MzA2MTQ2OSwiZXhwIjoxNjQzMDk3NDY5LCJub25jZSI6Ikh4TU43Y2RSbEhFWVlGQlBVNEJMUXY3STdKOE5ZMm92eDVRVUI3ZnVsWFkifQ.iBN9AUedfIo8AuXG3-YbJ8PB1xCR7avV49dBPGDTjkzhh0c_VwXLIJaAqotxXRQJFx7J8sRC_iq6ej4uvHbSnu563GmTazg0Sj1Vuf3uwa81YglJiSfpZNvmfCCMROJtk1ACB7dj54bMPWLXVlA-nqMqp6Q7WJtRGNFesQ9Cra__6T1awvHgq1rlg8hfCDeHpsx3K-D8yoUB1sDxLk3tdVVvJnWm01xM9cNn3Fr51nde6aoNaaijISHZoMkSEb7sNfOVtdqItG3WlHe8LNisGrwuWJF0kiqC_4gfCodF5PAe-3-dRuELGpD4s8hitIs1tGxm-qz6ZfyZWEdL5alfjQ
```

However, if you try to access the `/private` URL, you'll get a `403 Access denied` response. Looking at the
logs of the NGAC enforcer we can see:

```bash
$ kubectl logs -n ngac -l app=ngac-authz
2022/01/24 23:51:28  debug	Resolved by: jwt_bearer, Engineering [scope="pep"]
2022/01/24 23:51:28  debug	Resolved by: request_path, /private [scope="pep"]
2022/01/24 23:51:28  debug	Resolved by: method, GET [scope="pep"]
2022/01/24 23:51:28  debug	check(Engineering [GET] -> /private) = false [scope="pep"]
2022/01/24 23:51:28  debug	PC(http) [scope="pdp/audit"]
2022/01/24 23:51:28  debug	deny(access denied) [scope="pep"]
```

It has read the `group` claim value and checked against the [NGAC graph](graph.txt) if the group has
permissions on the requested URI, but access is denied.

Let's see what happens if the user is moved to the **Admins** group. To do so:

* Go to the management console in Auth0, and modify the user's **group** claim from Engineering to **Admins** (note that all
the values are case sensitive).
* In the browser, logout again by going to the `/logout` path to go back to the login screen to get a new token.
* Once that is done, login and access the `/private` endpoint again. The request should succeed and you should
  see something like:
  ```
  Welcome, Ignasi!
  Group: Admins
  Accessing: /private

  Authenticated with token:
  eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkRGRWVyODZqY2lRQTNfUVdETkE3MyJ9.eyJodHRwczovL3p0YS1kZW1vL2NvdW50cnkiOiJFUyIsImh0dHBzOi8venRhLWRlbW8vZ3JvdXAiOiJBZG1pbnMiLCJuaWNrbmFtZSI6ImlnbmFzaSt0ZXN0IiwibmFtZSI6IklnbmFzaSIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8wNDRjMjU1MDk4NDE2M2M5OTQ3N2NkM2QyYjY0NWViND9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmlnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIyLTAxLTI0VDIzOjQyOjQwLjQ0NloiLCJlbWFpbCI6ImlnbmFzaSt0ZXN0QHRldHJhdGUuaW8iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiaXNzIjoiaHR0cHM6Ly9uYWN4LWRtei5ldS5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NjFlN2Y0NWE3NmRjM2EwMDZhYWU1MGRiIiwiYXVkIjoiZHl5VzBtZTRMcThuM3RZMzBGYXR1RFFGWHB0WnZtNG8iLCJpYXQiOjE2NDMwNjc3NjEsImV4cCI6MTY0MzEwMzc2MSwibm9uY2UiOiJKa096aHNMS181V1JOMFJLSmUxb0d6V3UwVXNqOENCaWhJbkRNbm5zSzVJIn0.lyJ3Z70wj3VjqLOD5CF3yZqOA4kcyogEy_xRG82uMEYKSaQ3eXb71ZRRQtsZyosjycKrEr3-9HKHrLUlsV1NDyooNlFkMPls__Li0MkDg3cpokQ17m1V5B6NCcisN4arJIzYawCb9pbKqnOpOwKrNdt8Um7g2TDMD2ZrsuFpfBva_O_a6Z8JY3f6QDShSazKA51F8URatq4ZydxjxMGw96qa46X6DpKtDR2vFwDb78fu2RGehst_KKWFxzTU7mF5aF_7cVIqwvqxpWsglGIteZqS2B8JA1QSK2pRZOpNER7ciCwftoYx8wgxYJmfEaHfqV4YvmBvqNhV9FbSEwct3A
  ```

Checking the NGAC enforcer logs we'll see that now it is allowing access:

```bash
$ kubectl logs -n ngac -l app=ngac-authz
2022/01/24 23:50:33  debug	Resolved by: jwt_bearer, Admins [scope="pep"]
2022/01/24 23:50:33  debug	Resolved by: request_path, /private [scope="pep"]
2022/01/24 23:50:33  debug	Resolved by: method, GET [scope="pep"]
2022/01/24 23:50:33  debug	check(Admins [GET] -> /private) = true [scope="pep"]
2022/01/24 23:50:33  debug	PC(http) [scope="pdp/audit"]
2022/01/24 23:50:33  debug	  Admins-privileged-protected-/private ops=[GET] [scope="pdp/audit"]
2022/01/24 23:50:33  debug	allow() [scope="pep"]
```

## Environment cleanup

To cleanup the environment and uninstall the NGAC enforcers, you can use
the following command:

```bash
$ make uninstall
```
