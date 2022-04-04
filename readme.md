### Be Aware: Still WIP

Usecase:
you have a file named `.git-secrets.json` in your repository:

```
{
  "$schema": "../../schema/def/v1.json",
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
- [ ] Also store config values in .git-secrets.yaml
- [x] YAML Schema validation via JSON Schema
- [ ] More code documentation
- [x] Secret min requirements
- [ ] File watches / Daemon
- [ ] Add more examples
- [ ] CI/CD via Github Actions
- [ ] Brew Repository

Todo Beta Release
- [ ] Unit Testing
- [ ] More stable API
- [ ] Private Key encoding
- [ ] Git commit hook to scan for decoded secrets


 