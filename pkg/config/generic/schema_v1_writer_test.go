package config_generic

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

type GetSchemaFunc = func() V1Schema

func CreateGetSchemaFunc(t *testing.T, fs afero.Fs, configPath string) GetSchemaFunc {
	return func() V1Schema {
		fileBytes, errRead := afero.ReadFile(fs, configPath)
		assert.NoError(t, errRead)

		var Parsed V1Schema
		errParse := json.Unmarshal(fileBytes, &Parsed)
		assert.NoError(t, errParse)

		return Parsed
	}
}

func NewWrappedV1Writer(t *testing.T, inputFileName string) (writer *V1Writer, originalSchema V1Schema, getSchema GetSchemaFunc) {

	fs := afero.NewMemMapFs()
	configPath := ".git-secrets.json"

	// copy file from embed fs to afero fs -> .git-secrets.json
	wantedConfig, errRead := testFiles.ReadFile(fmt.Sprintf("test_fs/%s", inputFileName))

	var wantedConfigParsed V1Schema
	errParse := json.Unmarshal(wantedConfig, &wantedConfigParsed)
	assert.NoError(t, errParse)

	assert.NoError(t, errRead)
	assert.NoError(t, afero.WriteFile(fs, configPath, wantedConfig, 0664))

	return NewV1Writer(fs, wantedConfigParsed, configPath), wantedConfigParsed, CreateGetSchemaFunc(t, fs, configPath)

}

func TestNewV1Writer(t *testing.T) {
	newWriter := NewV1Writer(afero.NewMemMapFs(), V1Schema{}, "")
	assert.NotNil(t, newWriter)
	assert.IsType(t, &V1Writer{}, newWriter)
}

func TestV1Writer_AddContext(t *testing.T) {

	t.Run("add context if not exists", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileBlankDefault)
		assert.Nil(t, original.Context["prod"])
		assert.NoError(t, writer.AddContext("prod"))

		newCtx := getSchema().Context["prod"]

		assert.NotNil(t, newCtx)
		assert.Nil(t, newCtx.Secrets)
		assert.Nil(t, newCtx.Configs)
		assert.Nil(t, newCtx.DecryptSecret)
	})

	t.Run("should not add context if already exists", func(t *testing.T) {
		writer, _, getSchema := NewWrappedV1Writer(t, TestFileBlankTwoContexts)
		assert.NotNil(t, getSchema().Context["prod"].Secrets)
		assert.Error(t, writer.AddContext("prod"))
		assert.NotNil(t, getSchema().Context["prod"].Secrets)
	})

}

func TestV1Writer_AddFileToRender(t *testing.T) {

	t.Run("create render files map if missing", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileBlankDefault)
		assert.Nil(t, original.RenderFiles)
		assert.NoError(t, writer.AddFileToRender("env", "fileIn", "fileOut"))
		newSchema := getSchema()
		assert.NotNil(t, newSchema.RenderFiles)
		assert.NotNil(t, newSchema.RenderFiles["env"])
		assert.NotNil(t, newSchema.RenderFiles["env"].Files)
		assert.Len(t, newSchema.RenderFiles["env"].Files, 1)
	})

	t.Run("should not add existing files to the same target", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileRealWorld)

		assert.NotNil(t, original.RenderFiles["env"])
		envFiles := original.RenderFiles["env"]
		assert.Len(t, envFiles.Files, 1)
		envFile := original.RenderFiles["env"].Files[0]
		assert.NotNil(t, envFile)
		assert.Equal(t, "templates/.env.dist", envFile.FileIn)
		assert.Equal(t, "templates/.env", envFile.FileOut)

		assert.Error(t, writer.AddFileToRender("env", "templates/.env.dist", "templates/.env"))
		assert.Len(t, getSchema().RenderFiles["env"].Files, 1)
	})

	t.Run("should add files to a target", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileRealWorld)

		assert.NotNil(t, original.RenderFiles["env"])
		envFiles := original.RenderFiles["env"]
		assert.Len(t, envFiles.Files, 1)

		assert.NoError(t, writer.AddFileToRender("env", "file-in", "file-out"))

		newSchema := getSchema()
		envFiles = newSchema.RenderFiles["env"]
		assert.Len(t, envFiles.Files, 2)

		newFile := envFiles.Files[1]
		assert.NotNil(t, newFile)
		assert.Equal(t, "file-in", newFile.FileIn)
		assert.Equal(t, "file-out", newFile.FileOut)

	})

}

func TestV1Writer_SetConfig(t *testing.T) {

	t.Run("fail if context does not exists", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileRealWorld)
		assert.Nil(t, original.Context["missing"])
		assert.Error(t, writer.SetConfig("missing", "databaseHost", "testvalue", false))
		assert.Nil(t, getSchema().Context["missing"])
	})

	t.Run("initialize configs struct", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileBlankDefault)
		assert.Nil(t, original.Context["default"].Configs)
		assert.NoError(t, writer.SetConfig("default", "testKey", "testValue", false))
		assert.NotNil(t, getSchema().Context["default"].Configs)
		assert.Equal(t, "testValue", getSchema().Context["default"].Configs["testKey"])
	})

	t.Run("only overwrite when force is set to true", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileRealWorld)

		assert.Equal(t, "3306", original.Context["default"].Configs["databasePort"])
		assert.Error(t, writer.SetConfig("default", "databasePort", "3307", false))
		assert.Equal(t, "3306", getSchema().Context["default"].Configs["databasePort"])

		assert.NoError(t, writer.SetConfig("default", "databasePort", "3307", true))
		assert.Equal(t, "3307", getSchema().Context["default"].Configs["databasePort"])

	})

	t.Run("should not add entry if not defined in default context", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileRealWorld)

		assert.Equal(t, "", original.Context["default"].Configs["databaseName"])
		assert.Equal(t, "", original.Context["prod"].Configs["databaseName"])
		assert.Error(t, writer.SetConfig("prod", "databaseName", "work-cli", false))
		assert.Equal(t, "", getSchema().Context["prod"].Configs["databaseName"])

		assert.NoError(t, writer.SetConfig("default", "databaseName", "work-cli", false))
		assert.Equal(t, "work-cli", original.Context["default"].Configs["databaseName"])

		assert.NoError(t, writer.SetConfig("prod", "databaseName", "work-cli", false))
		assert.Equal(t, "work-cli", getSchema().Context["prod"].Configs["databaseName"])

	})

	t.Run("should write the config entry", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileRealWorld)
		assert.Equal(t, "", original.Context["default"].Configs["databaseName"])
		assert.NoError(t, writer.SetConfig("default", "databaseName", "work-cli", false))
		assert.Equal(t, "work-cli", getSchema().Context["default"].Configs["databaseName"])
	})

}

func TestV1Writer_SetSecret(t *testing.T) {

	t.Run("fail if context does not exists", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileRealWorld)
		assert.Nil(t, original.Context["missing"])
		assert.Error(t, writer.SetEncryptedSecret("missing", "databaseHost", "<encryptedValue>", false))
		assert.Nil(t, getSchema().Context["missing"])
	})

	t.Run("initialize secrets struct", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileConfigEntries)
		assert.Nil(t, original.Context["default"].Secrets)
		assert.NoError(t, writer.SetEncryptedSecret("default", "testKey", "<encryptedValue>", false))
		assert.NotNil(t, getSchema().Context["default"].Secrets)
		assert.Equal(t, "<encryptedValue>", getSchema().Context["default"].Secrets["testKey"])
	})

	t.Run("only overwrite when force is set to true", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileRealWorld)

		assert.Equal(t, "prPy40oRzdeFelmL5xVhbadEWNV9puR3/aWTY+gTYXOrT2bksi5GS9lCTKi66A3ePYa0hbwMqXadlDZw", original.Context["default"].Secrets["databasePassword"])
		assert.Error(t, writer.SetEncryptedSecret("default", "databasePassword", "<encryptedValue>", false))
		assert.Equal(t, "prPy40oRzdeFelmL5xVhbadEWNV9puR3/aWTY+gTYXOrT2bksi5GS9lCTKi66A3ePYa0hbwMqXadlDZw", original.Context["default"].Secrets["databasePassword"])

		assert.NoError(t, writer.SetEncryptedSecret("default", "databasePassword", "<encryptedValue>", true))
		assert.Equal(t, "<encryptedValue>", getSchema().Context["default"].Secrets["databasePassword"])

	})

	t.Run("should not add entry if not defined in default context", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileRealWorld)

		assert.Equal(t, "", original.Context["default"].Secrets["databaseRootPassword"])
		assert.Equal(t, "", original.Context["prod"].Secrets["databaseRootPassword"])
		assert.Error(t, writer.SetEncryptedSecret("prod", "databaseRootPassword", "<encryptedValue>", false))
		assert.Equal(t, "", getSchema().Context["prod"].Secrets["databaseRootPassword"])

		assert.NoError(t, writer.SetEncryptedSecret("default", "databaseRootPassword", "<encryptedValue>", false))
		assert.Equal(t, "<encryptedValue>", original.Context["default"].Secrets["databaseRootPassword"])

		assert.NoError(t, writer.SetEncryptedSecret("prod", "databaseRootPassword", "<encryptedValue1>", false))
		assert.Equal(t, "<encryptedValue1>", getSchema().Context["prod"].Secrets["databaseRootPassword"])

	})

	t.Run("should write the secret entry", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileRealWorld)
		assert.Equal(t, "", original.Context["default"].Secrets["databaseRootPassword"])
		assert.NoError(t, writer.SetEncryptedSecret("default", "databaseRootPassword", "<encryptedValue>", false))
		assert.Equal(t, "<encryptedValue>", getSchema().Context["default"].Secrets["databaseRootPassword"])

		assert.NoError(t, writer.SetEncryptedSecret("prod", "databaseRootPassword", "<encryptedValue1>", false))
		assert.Equal(t, "<encryptedValue1>", getSchema().Context["prod"].Secrets["databaseRootPassword"])
	})

}

func TestV1Writer_WriteConfig(t *testing.T) {

	t.Run("should fail if the resulting schema would be invalid", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileConfigEntries)

		// let's hack an invalid config
		assert.Equal(t, "", original.Context["prod"].Configs["missingKey"])
		writer.schema.Context["prod"].Configs["missingKey"] = "abc"

		assert.Error(t, writer.WriteConfig())

		assert.Equal(t, "", getSchema().Context["prod"].Configs["missingKey"])
	})

	t.Run("should write the config file", func(t *testing.T) {
		writer, original, getSchema := NewWrappedV1Writer(t, TestFileConfigEntries)
		assert.Equal(t, "3306", original.Context["default"].Configs["databasePort"])
		writer.schema.Context["default"].Configs["databasePort"] = "3308"
		assert.NoError(t, writer.WriteConfig())
		assert.Equal(t, "3308", getSchema().Context["default"].Configs["databasePort"])
	})

}
