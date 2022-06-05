## With Binary example

### Features Used
- Encoding / Decoding
- File Rendering

### How to run

````bash
# Configure the Encoder Secret used in this example
git-secrets set global-secret withbinaryexample --value XahcoQuae0wie3nooy0vuneiyaiy6phe

# Use the prod context
git-secrets render env -c prod && go run main.go

# Expected Output:
# .env written
# Environment Used: .env
# Database Host: my-prod-database.svc.local
# Database Port: 3306
# Database Name: git-secrets-demo
# Database Password: koocoo4pohKix8sei3eeve5areixeide


# Use the default context
git-secrets render env && go run main.go

# Expected Output:
# .env written
# Environment Used: .env
# Database Host: my-local-database.svc.local
# Database Port: 3306
# Database Name: git-secrets-demo
# Database Password: sooRahvow9eeXei5Eeph7ax9lee4AiG1
````