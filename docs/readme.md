## encryption and rendering engine for git repositories

![Tag](https://img.shields.io/github/v/release/benammann/git-secrets?label=release)
![Docker Image Version](https://img.shields.io/docker/v/benammann/git-secrets?label=image)
![Release Badge](https://github.com/benammann/git-secrets/actions/workflows/goreleaser.yml/badge.svg)
![Test Badge](https://github.com/benammann/git-secrets/actions/workflows/docker-release.yml/badge.svg)
![License](https://img.shields.io/github/license/benammann/git-secrets?123)

Git Secrets encrypts your passwords and configurations for multiple environments and allows you to check them into a git repository. Using the GoLang templating engine, you can then decrypt them and write them to env files or Kubernetes deployment files.

### How does it work

- For each Project / Context you can use a **Encoder Secret** which is stored at `~/.git-secrets.yaml`
- The **Encoder Secret** is used to encode your passwords which are then stored inside your git repositories `.git-secrets.json`
- The encrypted secrets are then decoded and rendered using Go Web Templates like Helm for example. (https://gowebexamples.com/templates/)
- Each project can have multiple contexts for example `default` and `prod`
- Every custom context inherits from the `default` context, so you don't have to define values twice
- You can use a different **Encoder Secret** in each context so the engineer can only access the secrets he should need

### Demo

![](./img/git-secrets-demo.gif)

### Examples

- Encoding / Decoding: [with-binary-example](examples/with-binary-example)
- Kubernetes Secrets: [render-kubernetes-secret](examples/render-kubernetes-secret)
- Github Actions [.github/workflows/docker-release.yml](.github/workflows/docker-release.yml)


### Installation

`Git-Secrets` is available on Linux, macOS and Windows platforms.

* Binaries for Linux, Windows and Mac are available as tarballs in the [release](https://github.com/benammann/git-secrets/releases) page.


* Via Curl for Linux and Mac (uses https://github.com/jpillora/installer)

  ```shell
  # without sudo
  curl https://i.jpillora.com/benammann/git-secrets! | bash
  
  # using sudo (if mv fails)
  curl https://i.jpillora.com/benammann/git-secrets!! | bash
  ```

* Via Homebrew for macOS or LinuxBrew for Linux

   ```shell
   brew install benammann/tap/git-secrets 
   ```

* Via a GO install

  ```shell
  # NOTE: The dev version will be in effect!
  go install github.com/benammann/git-secrets@latest
  ```

## Getting started

### Initialize the project
The configuration is made in a json file called `.git-secrets.json` you can also specify a custom path using `-f <path-to-custom-file>`

```bash
# Create a new global encoder secret (which you can later share with your team)
git secrets set global-secret mySecret --value $(pwgen -c 32 -n -s -y)

# Create a new .git-secrets.json
git secrets init

# Get the initial information of the config file
git secrets info

# Get the CLI's current version
git secrets version
```

### Encode a secret and add a config entry

Git-Secrets allows you to store encrypted `Secrets` and plain `Configs` both are stored in `.git-secrets.json`

```bash
# Encode a value (uses interactive input)
git secrets set secret databasePassword

# Write the value to a custom context
# Add Context: git secrets add context dev
git secrets set secret databasePassword -c dev

# Add a new config value
git secrets set config databaseHost db-host.svc.local

# Write the config value to a custom context
# Add Context: git secrets add context dev
git secrets set config databaseHost db-host.my-dev-db.svc -c dev
```

### Decode the secrets and get the config entry

```bash
# Get the decoded value
git secrets get secret databasePassword

# Get the value stored in databaseHost
git secrets get config databaseHost
```

### Create a `.env.dist` file

Git-Secrets allows you to render files using the `Secret` and `Config` values on the fly using gotemplates, just like Helm. For a syntax reference head over to https://gowebexamples.com/templates/

````text
DATABASE_HOST={{.Configs.databaseHost}}
DATABASE_PASSWORD={{.Secrets.databasePassword}}
````

You can have custom renderTargets to render files. For example `env` or `k8s`. You can than add multiple files to a renderTargets.

````bash
# always render empty.dist to .env
# uses the targetName: env
git secrets add file empty.dist .env -t env

# now execute the rendering process
# this renders the empty.dist file to .env and fills out all variables using the default context
# targetName: env
git secrets render env

# prints all available variables
git secrets render env --debug

# prints the rendered files to the console without actually writing the file
git secrets render env --dry-run

# renders the files using the prod context
git secrets render env -c prod
````

### Scan for plain secrets

`Git-Secrets` provides a simple command to scan for plain secrets in the project files.

![](./img/git-secrets-scan-demo.png)

````bash
# scan all files added to git
git secrets scan -a

# scan staged files only
git secrets scan

# hint: add -v to show all the scanned file names
````

You should use this command to setup a pre-commit git-hook in your project. You can use Husky (https://typicode.github.io/husky/#/) to automatically install and setup the hook.


### Custom Template Functions

Git Secrets extends the GoLang Templating engine by some useful functions

#### Base64Encode

The Base64Encode function takes the first argument and encodes it as Base64. This allows you to render Kubernetes Secrets

````yaml
# Created by git-secrets
apiVersion: v1
data:
  apiPassword: "{{ Base64Encode .Secrets.applicationAPassword }}"
kind: Secret
metadata:
  name: api-application-a
  namespace: {{.Configs.namespace}}
type: Opaque
````

#### GitConfig

GitConfig allows you to resolve git config values. For example if you want to render files individually to the developer

````text
GIT_NAME={{GitConfig "user.name"}}
GIT_EMAIL={{GitConfig "user.email"}}
````


### Using Docker

There is also a Docker Image available: `benammann/git-secrets`.

Since git-secrets normally depends on a global `.git-secrets.yaml` you need to use the `--secret` parameter to pass the encryption secret using cli.
You also need to mount the project's `.git-secrets.json` file using docker volume mounts.

````bash
# just execute the help command
docker run benamnann/git-secrets help

# get all the information about the .git-secrets.json file
docker run \
  # mount .git-secrets.json to /git-secrets/.git-secrets.json
  -v $PWD/.git-secrets.json:/git-secrets/.git-secrets.json \
  # use the official docker image
  benammann/git-secrets \
  # execute the info command
  info
  
docker run \
  # mount .git-secrets.json to /git-secrets/.git-secrets.json
  -v $PWD/.git-secrets.json:/git-secrets/.git-secrets.json \
  # use the official docker image
  benammann/git-secrets \
  # pass the encryption secret 'gitsecretspublic' including it's value from an local Environment variable to docker
  --secret gitsecretspublic=${SECRET_VALUE} \
  # decrypt the secret crToken
  get secret crToken 
````

## Documentation

### How the encryption is done

Git-Secrets uses AES-256 to encrypt / decrypt the secrets. Read more about it here [Advanced Encryption Standard](https://de.wikipedia.org/wiki/Advanced_Encryption_Standard).

The encryption key is stored outside your git repository and can be referenced using multiple methods

The implementation can be found here [engine_aes.go](pkg/encryption/engine_aes.go).

#### Named Secrets
Named secrets are stored in `~/.git-secrets.yaml` and have a name. You can than reference it using the `context.decryptSecret.fromName` key.

````
"decryptSecret": {
    "fromName": "withbinaryexample"
},
````

You can define a `decryptSecret` in each context to for example encrypt the production secrets using a different encryption key. This can be useful to not let your developers know the CI/CD Secrets.

The CLI provides multiple ways how to configure and manage your global secrets.
```bash
# Generate via pwgen and read from stdin
git secrets set global-secret mySecret --value $(pwgen -c 32 -n -s -y)

# Set manually using interactive input
git secrets set global-secret mySecret

# Get the written secret
git secrets get global-secret mySecret

# Get all global secret names
git secrets get global-secrets
```

#### Overwrite using CLI Args

In case you don't want to store the secrets globally and on the disk you can also use the following cli args to inject the secrets at runtime

```bash
# Uses the secret passed via --secret (insecure)
git secrets get secret mySecret --secret secretName=$(SECRET_VALUE) --secret secretName1=$(SECRET_VALUE_1)
```