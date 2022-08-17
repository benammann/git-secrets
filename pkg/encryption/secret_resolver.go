package encryption

import (
	"fmt"
	global_config "github.com/benammann/git-secrets/pkg/config/global"
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
	globalConfig        *global_config.GlobalConfigProvider
	overwrites          map[string]string
}

func (m *MergedSecretResolver) GetPlainSecret() (secret []byte, errResolve error) {

	if m.overwrites[m.requestedSecretName] != "" {
		return []byte(m.overwrites[m.requestedSecretName]), nil
	}

	secretValue := m.globalConfig.GetSecret(m.requestedSecretName)
	if secretValue == "" {
		return nil, fmt.Errorf("secret %s can not be found globally. Either pass --secret %s=$(MY_SECRET_NAME) or configure it using git secret global-secret", m.requestedSecretName, m.requestedSecretName)
	}
	return []byte(secretValue), nil

}

func NewMergedSecretResolver(requestedSecretName string, globalConfig *global_config.GlobalConfigProvider, overwrites map[string]string) SecretResolver {
	return &MergedSecretResolver{
		requestedSecretName: requestedSecretName,
		globalConfig:        globalConfig,
		overwrites:          overwrites,
	}
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
