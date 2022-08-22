package config_generic

import (
	"embed"
	"fmt"
	global_config "github.com/benammann/git-secrets/pkg/config/global"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed test_fs
var testFiles embed.FS

const GlobalSecretKey = "gitSecretsTest"
const GlobalSecretValue = "eeSaoghoh8oi9leed7hai4looK3jae1N"

const TestFileBlankDefault = "generic_repository_test-blank-default.json"
const TestFileBlankTwoContexts = "generic_repository_test-blank-two-contexts.json"
const TestFileBlankInvalidVersion = "generic_repository_test-invalid-version.json"
const TestFileConfigEntries = "generic_repository_test-config-entries.json"
const TestFileMissingEncryptionSecret = "generic_repository_test-missing-encryption-secret.json"

func createTestRepository(fileName string, selectedContextName string) (*Repository, error) {
	fileName = fmt.Sprintf("test_fs/%s", fileName)
	globalConfig := global_config.NewGlobalConfigProvider(global_config.NewMemoryStorageProvider())
	_ = globalConfig.SetSecret(GlobalSecretKey, GlobalSecretValue, false)
	mergeGlobalSecrets := make(map[string]string)
	repository, errParse := ParseRepositoryFromReadFileFs(testFiles, fileName, globalConfig, mergeGlobalSecrets)
	if errParse != nil {
		return nil, fmt.Errorf("could not parse: %s", errParse.Error())
	}
	_, errSetContext := repository.SetSelectedContext(selectedContextName)
	if errSetContext != nil {
		return nil, fmt.Errorf("could not set context: %s", errSetContext.Error())
	}
	return repository, nil
}

func initRepository(t *testing.T, fileName string, selectedContextName string) *Repository {
	repo, errParse := createTestRepository(fileName, selectedContextName)
	assert.NotNil(t, repo)
	assert.NoError(t, errParse)
	return repo
}

func TestNewRepository(t *testing.T) {

	t.Run("it should parse a valid repository", func(t *testing.T) {
		repo, errParse := createTestRepository(TestFileBlankDefault, "default")
		assert.NotNil(t, repo)
		assert.NoError(t, errParse)
	})

	t.Run("it should not parse a repository with an invalid version number", func(t *testing.T) {
		repo, errParse := createTestRepository(TestFileBlankInvalidVersion, "default")
		assert.Nil(t, repo)
		assert.Error(t, errParse)
	})

}

func TestRepository_GetConfigVersion(t *testing.T) {
	repo, errParse := createTestRepository(TestFileBlankDefault, "default")
	assert.NotNil(t, repo)
	assert.NoError(t, errParse)
	assert.Equal(t, 1, repo.GetConfigVersion())
}

func TestRepository_GetConfigWriter(t *testing.T) {
	repo, errParse := createTestRepository(TestFileBlankDefault, "default")
	assert.NotNil(t, repo)
	assert.NoError(t, errParse)
	assert.IsType(t, &V1Writer{}, repo.GetConfigWriter())
}

func TestRepository_IsDefault(t *testing.T) {
	repo, errParse := createTestRepository(TestFileBlankTwoContexts, "default")
	assert.NotNil(t, repo)
	assert.NoError(t, errParse)
	assert.True(t, repo.IsDefault())

	missingContext, errSet := repo.SetSelectedContext("does-not-exist")
	assert.Error(t, errSet)
	assert.Nil(t, missingContext)

	prodContext, errSet := repo.SetSelectedContext("prod")
	assert.NoError(t, errSet)
	assert.Equal(t, "prod", prodContext.Name)

}
