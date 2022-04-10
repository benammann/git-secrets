## encryption and rendering engine for git repositories

Git-Secrets allows you to add encrypted values to your git repository and decode them locally. The encrypted version is left in the upstream, the decoded version kept locally.

The main benefit is that you can render templates using the decoded values like `.env` file or event kubernetes deployment files. More details: https://github.com/benammann/git-secrets/tree/dev-beta/examples

Be aware that this project is still under development and the api may change.

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

# Result: .git-secrets.json written

# Get the initial information of the config file
git-secrets info
```

Now you need to replace the value of `contet.default.decryptSecret.fromName` by your global secret name. Example `mySecret`

```json
{
  "context": {
    "default": {
      "decryptSecret": {
        "fromName": "mySecret"
      }
    }
  }
}
```

### Encode a secret and add it to the config file
You can encode secrets using the `git-secrets encode`

```bash
# Encode a value
git-secrets encode "git-secrets rocks"

# Result: G2gt4a+L4lJdYMJBG8f2eMsvTGMTFxduDoUwzwz/gdcmWvsBo+3Um+l9wOal
```

Now you need to put the result into your `.git-secrets.json` file. Example:

```json
{
  "context": {
    "default": {
      "secrets": {
        "testSecret": "G2gt4a+L4lJdYMJBG8f2eMsvTGMTFxduDoUwzwz/gdcmWvsBo+3Um+l9wOal"
      }
    }
  }
}
```

Now you can get it's decoded value using the following command

```bash
# Decode a value
git-secrets decode testSecret

# Result: git-secrets rocks
```

### More documentation is added soon

Usecase:
you have a file named `.git-secrets.json` in your repository:

```json
{
  "$schema": "schema/def/v1.json",
  "version": 1,
  "context": {
    "default": {
      "decryptSecret": {
        "fromName": "git-secrets.example-default"
      },
      "secrets": {
        "applicationAPassword": "FTRldAR9SOt0/LuIFPbc1t5SHjq91I9XmaCL5Dg/AWJzQ/DY3DG5blpVTLH4hZYk4t1w+SRn5O4GhLiu",
        "applicationBPassword": "dCUK7Jfd5aB+WcI64AgX0/I7yT/OGMoUD0+uGgp5cs/smJAFvUWdBohNgmHg9KC4ExzWrt9beuCRorXI"
      }
    }
  },
  "renderFiles": {
    "default": {
      "files": [
        {
          "fileIn": "application-a/.env.dist",
          "fileOut": "application-a/.env"
        },
        {
          "fileIn": "application-b/.env.dist",
          "fileOut": "application-b/.env"
        }
      ]
    }
  }
}
```

Decode Secrets globally defined at `~/.git-secrets.yaml`
```
secrets:
  git-secrets:
     example-default: eid1chux0shuo5iegoomei2Uhohsai6k
```

- `git-secrets decode applicationAPassword`: Decodes the value of applicationAPassword
- `git-secrets render`: Renders the files using the decoded values
- `git-secrets encode <value>`: Encodes the value using git-secrets.example-default
- `git-secrets info -d` Shows the decoded secrets in the terminal
- `git-secrets init`: Create a new .git-secrets file

Check the examples folder for more details

Features:
- Encodes secrets allowing you to keep them into your git repositories
- Render .env files or kubernetes files and also kubernetes secrets (using the Base64 encode method) locally or in CI/CD

Todo Alpha Release
- [x] Global Secret Management via CLI
- [x] Also store config values in .git-secrets.yaml
- [x] YAML Schema validation via JSON Schema
- [ ] More code documentation
- [x] Secret min requirements
- [x] File watches / Daemon
- [ ] Add more examples
- [x] Github Actions
- [x] Brew Repository
- [ ] Modify .git-secrets.json (add encoded secrets)

Todo Beta Release
- [ ] Unit Testing
- [ ] More stable API
- [ ] Private Key encoding
- [ ] Git commit hook to scan for decoded secrets


 