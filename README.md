An environment variable plugin that sign repository and build information and export it a as JWS token.

_Please note this project requires Drone server version 1.6 or higher._

## Usage

First, generate a signing key :

```console
DRONE_PRIVATE_KEY="$(openssl genpkey -algorithm ed25519 -outform PEM)"
```

Download and run the plugin:

```console
$ docker run -d \
  --publish=3000:80 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --env=DRONE_PRIVATE_KEY="$DRONE_PRIVATE_KEY" \
  --restart=always \
  --name=drone-env-merge
```

Update your runner configuration to include the plugin address and the shared secret.

```text
DRONE_ENV_PLUGIN_ENDPOINT=http://drone-env-signed:3000
DRONE_ENV_PLUGIN_TOKEN=bea26a2221fd8090ea38720fc445eca6
```

Two environment variables will be added
to your pipelines, `DRONE_SIGNED_BUILD` and
`DRONE_SIGNED_REPO`. The JSON structure is [defined
here](https://github.com/drone/drone-go/blob/master/drone/types.go)
in the `Repo` and `Build` structures.
