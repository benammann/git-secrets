package config_generic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type ConfigMapWithOriginContextName map[string]*ConfigValueWithOriginContext

type ConfigValueWithOriginContext struct {
	ConfigValue       string
	OriginContextName string
}

func TestRepository_AddConfig(t *testing.T) {

	repo := initRepository(t, TestFileBlankTwoContexts, "default")

	defaultCtx := repo.GetContext("default")
	prodCtx := repo.GetContext("prod")

	assert.NotNil(t, defaultCtx)
	assert.NotNil(t, prodCtx)

	t.Run("should fail when adding key to child context only", func(t *testing.T) {
		assert.Error(t, repo.AddConfig(&Config{
			Name:          "testKey",
			OriginContext: prodCtx,
			Value:         "testValue",
		}))
	})

	t.Run("should not fail when adding new key to default context", func(t *testing.T) {
		assert.NoError(t, repo.AddConfig(&Config{
			Name:          "testKey",
			OriginContext: defaultCtx,
			Value:         "testValue",
		}))
	})

	t.Run("should not faile when adding key to child context if defined in default", func(t *testing.T) {
		assert.NoError(t, repo.AddConfig(&Config{
			Name:          "testKey",
			OriginContext: prodCtx,
			Value:         "testValue",
		}))
	})

}

func TestRepository_GetConfigMap(t *testing.T) {

	repo := initRepository(t, TestFileConfigEntries, "default")

	t.Run("correctly create the config map for default context", func(t *testing.T) {
		defaultCtxConfigMap := repo.GetConfigMap()
		assert.NotNil(t, defaultCtxConfigMap["databaseHost"])
		assert.NotNil(t, defaultCtxConfigMap["databasePort"])
		assert.Equal(t, "database.svc.local", defaultCtxConfigMap["databaseHost"])
		assert.Equal(t, "3306", defaultCtxConfigMap["databasePort"])
	})

	t.Run("correctly should inherit databasePort from default and use databaseHost from default context", func(t *testing.T) {
		prodCtx, errSetProd := repo.SetSelectedContext("prod")
		assert.NoError(t, errSetProd)
		assert.NotNil(t, prodCtx)
		assert.Equal(t, "prod", prodCtx.Name)

		prodCtxConfigMap := repo.GetConfigMap()
		assert.NotNil(t, prodCtxConfigMap["databaseHost"])
		assert.NotNil(t, prodCtxConfigMap["databasePort"])
		assert.Equal(t, "database.svc.cluster", prodCtxConfigMap["databaseHost"])
		assert.Equal(t, "3306", prodCtxConfigMap["databasePort"])
	})

}

func TestRepository_GetConfigsByContext(t *testing.T) {

	repo := initRepository(t, TestFileConfigEntries, "default")

	t.Run("default context has two entries defined", func(t *testing.T) {
		defaultConfigs := repo.GetConfigsByContext("default")
		assert.Len(t, defaultConfigs, 2)
	})

	t.Run("prod context has one entry defined", func(t *testing.T) {
		prodConfigs := repo.GetConfigsByContext("prod")
		assert.Len(t, prodConfigs, 1)
	})

}

func TestRepository_GetCurrentConfig(t *testing.T) {

	repo := initRepository(t, TestFileConfigEntries, "default")

	defaultDatabaseHost := repo.GetCurrentConfig("databaseHost")
	assert.NotNil(t, defaultDatabaseHost)
	assert.Equal(t, "database.svc.local", defaultDatabaseHost.Value)
	assert.Equal(t, "default", defaultDatabaseHost.OriginContext.Name)

	prodCtx, errSetProd := repo.SetSelectedContext("prod")
	assert.NoError(t, errSetProd)
	assert.NotNil(t, prodCtx)
	assert.Equal(t, "prod", prodCtx.Name)

	prodDatabaseHost := repo.GetCurrentConfig("databaseHost")
	assert.NotNil(t, prodDatabaseHost)
	assert.Equal(t, "database.svc.cluster", prodDatabaseHost.Value)
	assert.Equal(t, "prod", prodDatabaseHost.OriginContext.Name)

	prodDatabasePort := repo.GetCurrentConfig("databasePort")
	assert.NotNil(t, prodDatabasePort)
	assert.Equal(t, "3306", prodDatabasePort.Value)
	assert.Equal(t, "default", prodDatabasePort.OriginContext.Name)

	assert.Nil(t, repo.GetCurrentConfig("missingEntry"))

}

func TestRepository_GetCurrentConfigs(t *testing.T) {

	configsToMap := func(configs []*Config) ConfigMapWithOriginContextName {
		mapOut := make(ConfigMapWithOriginContextName)
		for _, configItem := range configs {
			mapOut[configItem.Name] = &ConfigValueWithOriginContext{
				ConfigValue:       configItem.Value,
				OriginContextName: configItem.OriginContext.Name,
			}
		}
		return mapOut
	}

	repo := initRepository(t, TestFileConfigEntries, "default")

	defaultConfigMap := configsToMap(repo.GetCurrentConfigs())

	assert.NotNil(t, defaultConfigMap["databaseHost"])
	assert.NotNil(t, defaultConfigMap["databasePort"])

	assert.Equal(t, "database.svc.local", defaultConfigMap["databaseHost"].ConfigValue)
	assert.Equal(t, "default", defaultConfigMap["databaseHost"].OriginContextName)
	assert.Equal(t, "3306", defaultConfigMap["databasePort"].ConfigValue)
	assert.Equal(t, "default", defaultConfigMap["databasePort"].OriginContextName)

	_, errCtx := repo.SetSelectedContext("prod")
	assert.NoError(t, errCtx)

	prodConfigMap := configsToMap(repo.GetCurrentConfigs())
	assert.NotNil(t, prodConfigMap["databaseHost"])
	assert.NotNil(t, prodConfigMap["databasePort"])

	assert.Equal(t, "database.svc.cluster", prodConfigMap["databaseHost"].ConfigValue)
	assert.Equal(t, "prod", prodConfigMap["databaseHost"].OriginContextName)
	assert.Equal(t, "3306", prodConfigMap["databasePort"].ConfigValue)
	assert.Equal(t, "default", prodConfigMap["databasePort"].OriginContextName)

}
