# vaulter-white
A tool to pass Vault secrets to other processes via environment variables.

[![Build Status](https://travis-ci.org/cosmonawt/vaulter-white.svg?branch=master)](https://travis-ci.org/cosmonawt/vaulter-white)
[![Go Report Card](https://goreportcard.com/badge/github.com/cosmonawt/vaulter-white)](https://goreportcard.com/report/github.com/cosmonawt/vaulter-white)

## About
vaulter-white reads secrets from [Vault](https://vaultproject.io) and passes them into a newly spawned process using environment variables.
It is particularly useful in containerized applications. For example it can be set as the `ENTRYPOINT` in a Docker container to retrieve production Keys for your App.
After loading the secrets into the environment (while also including all existing variables) it will replace itself with a freshly spawned instance of the configurable process.

_Note:_ At the moment only [AppRole](https://www.vaultproject.io/docs/auth/approle.html) authentication is supported.

## Configuration
vaulter-white is configured via vaulter-white.yaml which will be read from the current directory if not specified otherwise using the `-c` flag.

_Example:_
```yaml
command: ["bash", "-c", "env"]            # Specifies the command to run after loading the secrets.
host: http://vault.rocks:8200             # Host of Vault server.
roleId: myAppRole                         # RoleID and SecretID for AppRole Authentication in Vault.
secretId: mySuperSecretId
secretIdEnv: SECRET_ID                    # The name of an environment variable storing the secretId, if not specified above.
secretMount: /secret/appConfig/           # secretMount contains the path to the secret backend holding your keys in Vault.
secrets:                                  # secrets is a collection of environment variable name overrides for each key.
  awsConfig:
    region: AWS_REGION
    access_key_id: AWS_KEY_ID
    secret_access_key: AWS_SECRET_KEY
  googleAPI:
    apiKey: GOOGLE_API_KEY

```

- `command` is optional and can be passed as command line argument as well (for example: `vaulter-white -c config.yaml bash -c env`).
- `secretId` will be read from environment variables (either at `secretIdEnv` as configured or at `VAULT_SECRET_ID`) if not configured. This makes it easy to include vaulter-white in Docker images that are built by CI.
- `secrets` is optional as well. Any keys not listed there will be exported as `SECRETNAME_KEY=value`.

_Note:_ Secret values should always store flat data types and no marshaled data (e.g JSON Objects). Values that are not strings will be exported as JSON.

## Run

To pass the configuration use the `-c` flag: `vaulter-white -c configuration.yaml`
If no command was specified in the configuration it should be passed as a commandline argument: `vaulter-white -c config.yaml bash -c env`
