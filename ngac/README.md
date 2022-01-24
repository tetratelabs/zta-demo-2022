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
