package encryption

import (
	"fmt"
	"github.com/benammann/git-secrets/pkg/cli_config"
	"github.com/spf13/viper"
	"os"
)

type SecretResolver interface {
	GetPlainSecret() (secret []byte, errResolve error)
}

type FromPlainSecretResolver struct {
	SecretResolver
	PlainSecret string
}

func NewPlainSecretResolver(plainSecret string) SecretResolver {
	return &FromPlainSecretResolver{
		PlainSecret: plainSecret,
	}
}

func(rs *FromPlainSecretResolver) GetPlainSecret() (secret []byte, errResolve error) {
	if rs.PlainSecret == "" {
		return nil, fmt.Errorf("got empty secret")
	}
	return []byte(rs.PlainSecret), nil
}

type FromEnvSecretResolver struct {
	SecretResolver
	envName string
}

func NewEnvSecretResolver(envName string) SecretResolver {
	return &FromEnvSecretResolver{
		envName: envName,
	}
}

func(rs *FromEnvSecretResolver) GetPlainSecret() (secret []byte, errResolve error) {
	envValue := os.Getenv(rs.envName)
	if envValue == "" {
		return nil, fmt.Errorf("env variable %s is empty", rs.envName)
	}
	return []byte(envValue), nil
}

type FromNameSecretResolver struct {
	SecretResolver
	secretName string
}

func NewNameSecretResolver(secretName string) SecretResolver {
	return &FromNameSecretResolver{
		secretName: secretName,
	}
}

func(rs *FromNameSecretResolver) GetPlainSecret() (secret []byte, errResolve error) {
	configKey := cli_config.NamedSecret(rs.secretName)
	configValue := viper.GetString(configKey)
	if configValue == "" {
		return nil, fmt.Errorf("no secret configured at %s", configKey)
	}
	return []byte(configValue), nil
}
