package cli_config

import "fmt"

const (
	Secrets = "secrets"
)

func NamedSecret(secretName string) string {
	return fmt.Sprintf("%s.%s", Secrets, secretName)
}