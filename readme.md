## encryption and rendering engine for git repositories

![Release Badge](https://github.com/benammann/git-secrets/actions/workflows/release.yml/badge.svg)
![Test Badge](https://github.com/benammann/git-secrets/actions/workflows/test.yml/badge.svg)

Git-Secrets allows you to add encrypted values to your git repository and decode them locally. The encrypted version is left in the upstream, the decoded version kept locally.

The main benefit is that you can render templates using the decoded values like `.env` file or event kubernetes deployment files. More details: https://github.com/benammann/git-secrets/tree/dev-beta/examples

Be aware that this project is still under development and the api may change.

### How does it work

- For each Project / Context you can use a **Encoder Secret** which is stored at `~/.git-secrets.yaml`
- The **Encoder Secret** is used to encode your passwords which are then stored inside your git repositories `.git-secrets.json`
- The encrypted secrets are then decoded and rendered using Go Web Templates like Helm for example. (https://gowebexamples.com/templates/)
- Each project can have multiple contexts for example `default` and `prod`
- Every custom context inherits from the `default` context, so you don't have to define values twice
- You can use a different **Encoder Secret** in each context so the engineer can only access the secrets he should need

### Examples

- Encoding / Decoding: [with-binary-example](examples/with-binary-example)
- Kubernetes Secrets: [render-kubernetes-secret](examples/render-kubernetes-secret)

### Installation

via Homebrew / Linuxbrew
```
brew install benammann/tap/git-secrets 
```

or just head over to the [Releases](https://github.com/benammann/git-secrets/releases) page and download the prebuilt binary manually

## Getting started

### Configure the global encoder secret
First, you need to create or configure a global secret. For this example our secret is called `mySecret`

Hint: Use a tool like `pwgen` to securely generate secrets locally. (Install via `brew install pwgen`)

Generate a global secret and set it to `mySecret`
```bash
# Generate via pwgen 
git-secrets global-secret mySecret $(pwgen -c 32 -n -s -y)

# Set manually
git-secrets global-secret mySecret <my-secret-here>

# Get the written secret
git-secrets global-secret mySecret

# Get all global secret names
git-secrets global-secret
```

### Initialize the project
The configuration is made in a json file called `.git-secrets.json` you can also specify a custom path using `-f <path-to-custom-file>`

```bash
# Create a new .git-secrets.json
git-secrets init

# .git-secrets.json written
# Info: git-secrets info -d
# Add Context: git-secrets add-context <contextName>
# Add Secret: git-secrets encode --write secretName

# Get the initial information of the config file
git-secrets info
```

### Encode a secret and add it to the config file
You can encode secrets using the `git-secrets encode`

```bash
# Encode a value (uses interactive input)
git-secrets set secret myAwesomeSecret

# ? Value to encode *****************
# Secret myAwesomeSecret written to .git-secrets.json
# Get the decoded value: git-secrets decode myAwesomeSecret
```

Now you can get it's decoded value using the following command

```bash
# Decode a value
git-secrets get secret myAwesomeSecret

# Result: Git Secrets Rocks
```