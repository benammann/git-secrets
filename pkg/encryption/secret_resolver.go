package encryption

import (
	"fmt"
	"github.com/benammann/git-secrets/pkg/config/cli"
	"github.com/spf13/viper"
	"os"
)

type SecretResolver interface {
	GetPlainSecret() (secret []byte, errResolve error)
}

type ConfigResolver interface {
	GetString(key string) string
}

type MergedSecretResolver struct {
	requestedSecretName string
	configResolver      ConfigResolver
	overwrites          map[string]string
}

func (m *MergedSecretResolver) GetPlainSecret() (secret []byte, errResolve error) {

	if m.overwrites[m.requestedSecretName] != "" {
		return []byte(m.overwrites[m.requestedSecretName]), nil
	}

	configKey := cli_config.NamedSecret(m.requestedSecretName)
	configValue := m.configResolver.GetString(configKey)
	if configValue == "" {
		return nil, fmt.Errorf("secret %s can not be found globally. Either pass --secret %s=$(MY_SECRET_NAME) or configure it using git secret global-secret", configKey, configKey)
	}
	return []byte(configValue), nil

}

func NewMergedSecretResolver(requestedSecretName string, configResolver ConfigResolver, overwrites map[string]string) SecretResolver {
	return &MergedSecretResolver{
		requestedSecretName: requestedSecretName,
		configResolver:      configResolver,
		overwrites:          overwrites,
	}
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

func (rs *FromPlainSecretResolver) GetPlainSecret() (secret []byte, errResolve error) {
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

func (rs *FromEnvSecretResolver) GetPlainSecret() (secret []byte, errResolve error) {
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

func (rs *FromNameSecretResolver) GetPlainSecret() (secret []byte, errResolve error) {
	configKey := cli_config.NamedSecret(rs.secretName)
	configValue := viper.GetString(configKey)
	if configValue == "" {
		return nil, fmt.Errorf("no secret configured at %s", configKey)
	}
	return []byte(configValue), nil
}
