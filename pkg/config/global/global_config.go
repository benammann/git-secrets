package global_config

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const SecretKeyPrefix = "secrets"

type GlobalConfigProvider struct {
	storageProvider StorageProvider
}

func NewGlobalConfigProvider(storageProvider StorageProvider) *GlobalConfigProvider {
	return &GlobalConfigProvider{
		storageProvider: storageProvider,
	}
}

func (g *GlobalConfigProvider) GetSecret(secretKey string) (value string) {
	return g.storageProvider.GetString(g.secretConfigKey(secretKey))
}

func (g *GlobalConfigProvider) SetSecret(secretKey string, secretValue string, force bool) error {

	errValidate := g.validateSecret(secretKey, secretValue)
	if errValidate != nil {
		return errValidate
	}

	configKey := g.secretConfigKey(secretKey)

	exists := g.GetSecret(secretKey) != ""
	if exists && force == false {
		return fmt.Errorf("secret %s already exists. use --force to overwrite", configKey)
	}

	g.storageProvider.Set(configKey, secretValue)

	return g.storageProvider.WriteConfig()
}

func (g *GlobalConfigProvider) GetSecretKeys() []string {
	var secretKeys []string
	for _, key := range g.storageProvider.AllKeys() {
		secretPrefix := fmt.Sprintf("%s.", SecretKeyPrefix)
		if strings.HasPrefix(key, secretPrefix) {
			secretKeys = append(secretKeys, strings.Replace(key, secretPrefix, "", 1))
		}
	}
	sort.Strings(secretKeys)
	return secretKeys
}

func (g *GlobalConfigProvider) secretConfigKey(secretKey string) string {
	return fmt.Sprintf("%s.%s", SecretKeyPrefix, strings.ToLower(secretKey))
}

func (g *GlobalConfigProvider) validateSecret(secretKey string, secretValue string) error {

	if !regexp.MustCompile(`^[A-Za-z1-9]+$`).MatchString(secretKey) {
		return fmt.Errorf("invalid key: only alphanumeric letters allowed [A-Za-z1-9] allowed")
	}

	if isInvalid := validateAESSecret(secretValue); isInvalid != nil {
		return fmt.Errorf("invalid value: %s", isInvalid.Error())
	}
	return nil
}

func validateAESSecret(plainSecret string) error {
	k := len([]byte(plainSecret))
	switch k {
	default:
		return fmt.Errorf("only key size of either 16, 24, or 32 bytes allowed")
	case 16, 24, 32:
		return nil
	}
}
