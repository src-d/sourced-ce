# How to enable OAuth

Following this guide, you will let the users to register and login using Google and Github OAUTH.

- active `docker-compose.yml` (`~/.sourced/compose-files/__active__/docker-compose.yml`) should contain:
    ```yml
    version: '3.4'
    services:
    sourced-ui:
      environment:
        OAUTH_ENABLED_PROVIDERS: ${OAUTH_ENABLED_PROVIDERS}           # Comma separated list of available OAuth providers   (eg: github,google)
        OAUTH_REGISTRATION_ROLE: ${OAUTH_REGISTRATION_ROLE}           # The role for the newly registered users using OAuth   'Admin'/'Alpha'/'Gamma'
        OAUTH_GOOGLE_CONSUMER_KEY: ${OAUTH_GOOGLE_CONSUMER_KEY}       # OAuth provider consumer key (aka client_id)
        OAUTH_GOOGLE_CONSUMER_SECRET: ${OAUTH_GOOGLE_CONSUMER_SECRET} # OAuth provider consumer secret (aka client_secret)
        OAUTH_GITHUB_CONSUMER_KEY: ${OAUTH_GITHUB_CONSUMER_KEY}       # OAuth provider consumer key (aka client_id)
        OAUTH_GITHUB_CONSUMER_SECRET: ${OAUTH_GITHUB_CONSUMER_SECRET} # OAuth provider consumer secret (aka client_secret)
    ```

In order to avoid some OAuth issues, make your instance of source{d} CE accessible from the Internet:
- consider this example: `http://live.sourced.ce`
    - `PROTOCOL` = `http`
    - `HOST` = `live.sourced.ce`
- you can use `ngrok` if running source{d} locally:
    ```shell
    $ ngrok http 8088
    ```

## For Google OAuth

1. Configure your "OAuth consent screen"
    - https://console.developers.google.com/apis/credentials/consent
    - set up:
        - "Application name"
        - "Authorized domains": `${HOST}`
1. Create a new "Credential"
    - https://console.developers.google.com/apis/credentials
    - create a new "Create OAuth client ID" for "Wab application"
    - set up:
        - "Application name"
        - "Restrictions":
            - Authorized JavaScript origins: `${PROTOCOL}://${HOST}`
            - Authorized redirect URIs: `${PROTOCOL}://${HOST}/oauth-authorized/google`


## For GitHub OAuth

1. Create a new OAuth application
    - https://github.com/settings/applications/new
    - set up:
        - "Application name"
        - "Homepage URL": `${PROTOCOL}://${HOST}`
        - "Authorization callback URL": `${PROTOCOL}://${HOST}/oauth-authorized/github`



## Run source{d}

```shell
$ export OAUTH_ENABLED_PROVIDERS=github,google
$ export OAUTH_REGISTRATION_ROLE=Gamma
$ OAUTH_GOOGLE_CONSUMER_KEY=<google client_id> \
  OAUTH_GOOGLE_CONSUMER_SECRET=<google client_secret> \
  OAUTH_GITHUB_CONSUMER_KEY=<github client_id> \
  OAUTH_GITHUB_CONSUMER_SECRET=<github client_secret> \
  sourced init local .
```

Every new user will be registered under `Gamma` role, so you might want to assign them a new role in order to let them to view `gitbase` database.

You can assign roles, or convert `Gamma` user into a an `Admin` one, following these steps:

---

_**[Disclaimer] [Bug]** Listing users is buggy and it's behavior depends on how you ran `sourced` for the first time:_

---

- If the first time you ran `sourced` was with no `OAUTH_ENABLED_PROVIDERS`, then you can only list users using `admin` user, so do the following:
    1. start `sourced` with `OAUTH_ENABLED_PROVIDERS=` (it will stop your current instance of `sourced`),
    1. login as `admin` (password: `admin`), and manage the users at you will.
    1. start `sourced` again, with the regular `OAUTH_ENABLED_PROVIDERS=github,google`,

- If the first time you ran `sourced` was with already valid `OAUTH_ENABLED_PROVIDERS`, do the following:
    1. start `sourced` with `OAUTH_REGISTRATION_ROLE=Admin` (it will stop your current instance of `sourced`),
    1. register as a new user, who will be created with `Admin` role, and who will be able to list users `/users/list`,
    1. start `sourced` again, with the regular `OAUTH_REGISTRATION_ROLE=Gamma` (it will stop the privileged instance of `sourced` from the previous step),
    1. login with the user that you created in the second step, and use it to manage the users at you will.
