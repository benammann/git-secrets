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
		writer, _, getSchema := NewWrappedV1Writer(t, TestFileBlankDefault)
		assert.Nil(t, getSchema().Context["prod"])
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

	t.Run("create files entry if missing", func(t *testing.T) {

	})

	t.Run("should not add existing files to the same target", func(t *testing.T) {

	})

	t.Run("should add files to a target", func(t *testing.T) {

	})

}

func TestV1Writer_SetConfig(t *testing.T) {

}

func TestV1Writer_SetSecret(t *testing.T) {

}

func TestV1Writer_WriteConfig(t *testing.T) {

}
