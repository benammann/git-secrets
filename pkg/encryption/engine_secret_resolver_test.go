package encryption

import (
	global_config "github.com/benammann/git-secrets/pkg/config/global"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFromEnvSecretResolver_GetPlainSecret(t *testing.T) {

	globalConfig := global_config.NewGlobalConfigProvider(global_config.NewMemoryStorageProvider())
	_ = globalConfig.SetSecret("overwritten", "riz9ohg9IefeeG8sha0quoa6it6uan6b", false)
	_ = globalConfig.SetSecret("original", "Ohqu7lahn4AiQu3reecoo1ausoo7aiy0", false)
	mergeGlobalSecrets := make(map[string]string)
	mergeGlobalSecrets["overwritten"] = "iepheam7aech9Wah5ahng5aix5Thumai"

	t.Run("should resolve the overwritten value", func(t *testing.T) {
		sr := NewMergedSecretResolver("overwritten", globalConfig, mergeGlobalSecrets)
		value, err := sr.GetPlainSecret()
		assert.NoError(t, err)
		assert.Equal(t, []byte("iepheam7aech9Wah5ahng5aix5Thumai"), value)
	})
	t.Run("should resolve the original value", func(t *testing.T) {
		sr := NewMergedSecretResolver("original", globalConfig, mergeGlobalSecrets)
		value, err := sr.GetPlainSecret()
		assert.NoError(t, err)
		assert.Equal(t, []byte("Ohqu7lahn4AiQu3reecoo1ausoo7aiy0"), value)
	})
	t.Run("should fail if secret does not exists", func(t *testing.T) {
		sr := NewMergedSecretResolver("missing", globalConfig, mergeGlobalSecrets)
		_, err := sr.GetPlainSecret()
		assert.Error(t, err)
	})

}

func TestMergedSecretResolver_GetPlainSecret(t *testing.T) {
	t.Run("should return env value", func(t *testing.T) {
		assert.NoError(t, os.Setenv("ENV_NAME", "value"))
		sr := NewEnvSecretResolver("ENV_NAME")
		value, err := sr.GetPlainSecret()
		assert.NoError(t, err)
		assert.Equal(t, []byte("value"), value)
	})
	t.Run("should fail if env is not set", func(t *testing.T) {
		sr := NewEnvSecretResolver("MISSING")
		_, err := sr.GetPlainSecret()
		assert.Error(t, err)
	})
}

func TestNewEnvSecretResolver(t *testing.T) {
	sr := NewEnvSecretResolver("ENV_NAME")
	assert.NotNil(t, sr)
}

func TestNewMergedSecretResolver(t *testing.T) {
	sr := NewMergedSecretResolver("secretName", nil, nil)
	assert.NotNil(t, sr)
}
