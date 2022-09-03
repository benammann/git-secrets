package config_generic

import (
	"encoding/json"
	"fmt"
	global_config "github.com/benammann/git-secrets/pkg/config/global"
	"github.com/benammann/git-secrets/pkg/encryption"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ParseAsSchemaV1(t *testing.T, fileName TestFile) V1Schema {
	fs := afero.FromIOFS{
		FS: testFiles,
	}
	fileBytes, errRead := afero.ReadFile(fs, fmt.Sprintf("test_fs/schema/v1/%s", fileName))
	assert.NoError(t, errRead)

	var Parsed V1Schema
	errParse := json.Unmarshal(fileBytes, &Parsed)
	assert.NoError(t, errParse)
	return Parsed

}

func TestIsSchemaV1(t *testing.T) {
	assert.True(t, IsSchemaV1(1))
	assert.False(t, IsSchemaV1(2))
	assert.False(t, IsSchemaV1(0))
}

func TestParseSchemaV1(t *testing.T) {

}

func TestV1Schema_validateSchemaV1(t *testing.T) {
	t.Run("fail on unsupported versions", func(t *testing.T) {
		parsed := ParseAsSchemaV1(t, "v2.json")
		assert.Equal(t, parsed.Version, 2)
		assert.Error(t, parsed.validateSchemaV1())
	})
	t.Run("fail if default ctx is missing", func(t *testing.T) {
		parsed := ParseAsSchemaV1(t, "no-default-ctx.json")
		assert.Error(t, parsed.validateSchemaV1())
	})
	t.Run("fail is multiple decrypt secret methods passed", func(t *testing.T) {
		parsed := ParseAsSchemaV1(t, "many-decrypt-secrets.json")
		assert.NotEqual(t, "", parsed.Context["default"].DecryptSecret.FromEnv)
		assert.NotEqual(t, "", parsed.Context["default"].DecryptSecret.FromName)
		assert.Error(t, parsed.validateSchemaV1())
	})
	t.Run("fail if no decrypt method on default ctx", func(t *testing.T) {
		parsed := ParseAsSchemaV1(t, "no-decrypt-secret-on-default.json")
		assert.Equal(t, "", parsed.Context["default"].DecryptSecret.FromEnv)
		assert.Equal(t, "", parsed.Context["default"].DecryptSecret.FromName)
		assert.Error(t, parsed.validateSchemaV1())
	})
	t.Run("fail if secret is defined in child but not in default context", func(t *testing.T) {
		parsed := ParseAsSchemaV1(t, "secret-missing-in-default.json")
		assert.Equal(t, "", parsed.Context["default"].Secrets["test"])
		assert.NotEqual(t, "", parsed.Context["prod"].Secrets["test"])
		assert.Error(t, parsed.validateSchemaV1())
	})
	t.Run("fail if config is defined in child but not in default context", func(t *testing.T) {
		parsed := ParseAsSchemaV1(t, "config-missing-in-default.json")
		assert.Equal(t, "", parsed.Context["default"].Configs["test"])
		assert.NotEqual(t, "", parsed.Context["prod"].Configs["test"])
		assert.Error(t, parsed.validateSchemaV1())
	})
	t.Run("do not fail if schema is valid", func(t *testing.T) {
		parsed := ParseAsSchemaV1(t, "real-world.json")
		assert.NoError(t, parsed.validateSchemaV1())
	})
}

func Test_getSecretResolverV1(t *testing.T) {

	globalConfig := global_config.NewGlobalConfigProvider(global_config.NewMemoryStorageProvider())
	_ = globalConfig.SetSecret(GlobalSecretKey, GlobalSecretValue, false)
	mergeGlobalSecrets := make(map[string]string)

	t.Run("return secret resolver from default context if no config given", func(t *testing.T) {
		repo := initRepository(t, TestFileBlankDefault, "default")
		defaultCtx := repo.GetContext("default")
		assert.NotNil(t, defaultCtx)
		assert.NotNil(t, defaultCtx.SecretResolver)
		assert.Equal(t, defaultCtx.SecretResolver, getSecretResolverV1(nil, defaultCtx, globalConfig, mergeGlobalSecrets))
	})
	t.Run("return from env secret resolver", func(t *testing.T) {
		repo := initRepository(t, TestFileBlankDefaultFromEnv, "default")
		defaultCtx := repo.GetContext("default")
		assert.NotNil(t, defaultCtx)
		assert.NotNil(t, defaultCtx.SecretResolver)
		assert.IsType(t, &encryption.FromEnvSecretResolver{}, defaultCtx.SecretResolver)
	})
	t.Run("return from name secret resolver", func(t *testing.T) {
		repo := initRepository(t, TestFileBlankDefault, "default")
		defaultCtx := repo.GetContext("default")
		assert.NotNil(t, defaultCtx)
		assert.NotNil(t, defaultCtx.SecretResolver)
		assert.IsType(t, &encryption.MergedSecretResolver{}, defaultCtx.SecretResolver)
	})
}
