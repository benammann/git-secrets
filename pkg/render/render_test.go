package render

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	global_config "github.com/benammann/git-secrets/pkg/config/global"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

const GlobalSecretKey = "gitSecretsTest"
const GlobalSecretValue = "eeSaoghoh8oi9leed7hai4looK3jae1N"

const FileRenderTestDefault = "render-test-default.json"

//go:embed test_fs
var testFiles embed.FS
var testFileAfero = &afero.FromIOFS{
	FS: testFiles,
}

type TestRenderFunctions struct {
	Base64Encode string `json:"base64Encode"`
}

type TestRenderContext struct {
	RenderContext *RenderingContext `json:"context"`
	TestFunctions *TestRenderFunctions `json:"functions"`
}

func createTestRepository(fileName string, selectedContextName string) (*config_generic.Repository, error) {
	fileName = fmt.Sprintf("test_fs/%s", fileName)
	globalConfig := global_config.NewGlobalConfigProvider(global_config.NewMemoryStorageProvider())
	_ = globalConfig.SetSecret(GlobalSecretKey, GlobalSecretValue, false)
	mergeGlobalSecrets := make(map[string]string)
	repository, errParse := config_generic.ParseRepository(afero.FromIOFS{
		FS: testFiles,
	}, fileName, globalConfig, mergeGlobalSecrets)
	if errParse != nil {
		return nil, fmt.Errorf("could not parse: %s", errParse.Error())
	}
	_, errSetContext := repository.SetSelectedContext(selectedContextName)
	if errSetContext != nil {
		return nil, fmt.Errorf("could not set context: %s", errSetContext.Error())
	}
	return repository, nil
}

func initRepository(t *testing.T, fileName string, selectedContextName string) (*config_generic.Repository, *RenderingEngine) {
	repo, errParse := createTestRepository(fileName, selectedContextName)
	assert.NotNil(t, repo)
	assert.NoError(t, errParse)
	return repo, NewRenderingEngine(repo, testFileAfero, testFileAfero)
}

func TestNewRenderingEngine(t *testing.T) {
	_, engine := initRepository(t, FileRenderTestDefault, "default")
	assert.NotNil(t, engine)
}

func TestRenderingEngine_CreateRenderingContext(t *testing.T) {
	repo, engine := initRepository(t, FileRenderTestDefault, "default")
	assert.NotNil(t, engine)

	file := &config_generic.FileToRender{
		FileIn: "fileIn",
		FileOut: "fileOut",
	}

	dbPassword := repo.GetCurrentSecret("databasePassword")
	dbPasswordVal, errDecode := dbPassword.Decode()
	assert.NoError(t, errDecode)

	ctx, err := engine.CreateRenderingContext(file)
	assert.NoError(t, err)
	assert.Equal(t, "default", ctx.ContextName)
	assert.Equal(t, dbPasswordVal, ctx.Secrets["databasePassword"])
	assert.Equal(t, "3306", ctx.Configs["databasePort"])
}

func TestRenderingEngine_ExecuteTemplate(t *testing.T) {
	_, engine := initRepository(t, FileRenderTestDefault, "default")
	fileToRender := &config_generic.FileToRender{
		FileIn: "test_fs/templates/render-context.json",
	}

	var bytesOut bytes.Buffer
	usedContext, errExecute := engine.ExecuteTemplate(fileToRender, &bytesOut)
	assert.NoError(t, errExecute)
	assert.NotNil(t, usedContext)

	var renderedFileDecoded TestRenderContext
	errDecode := json.Unmarshal(bytesOut.Bytes(), &renderedFileDecoded)
	assert.NoError(t, errDecode)
	assert.NotNil(t, renderedFileDecoded.RenderContext)

	renderContext := renderedFileDecoded.RenderContext
	assert.Equal(t, "default", usedContext.ContextName)
	assert.Equal(t, usedContext.ContextName, renderContext.ContextName)
	assert.Equal(t, usedContext.File.FileIn, renderContext.File.FileIn)
	assert.Equal(t, usedContext.File.FileOut, renderContext.File.FileOut)
	assert.Equal(t, "3306", usedContext.Configs["databasePort"])
	assert.Equal(t, usedContext.Secrets["databasePort"], renderContext.Secrets["databasePort"])
	assert.Equal(t, "em8toheGhieh0Thu1ahz9Lou2ucheeh6", usedContext.Secrets["databasePassword"])
	assert.Equal(t, usedContext.Secrets["databasePassword"], renderContext.Secrets["databasePassword"])

}

func TestRenderingEngine_RenderFile(t *testing.T) {
	_, engine := initRepository(t, FileRenderTestDefault, "default")
	fileToRender := &config_generic.FileToRender{
		FileIn: "test_fs/templates/render-context.json",
	}
	_, fileContents, err := engine.RenderFile(fileToRender)
	assert.NoError(t, err)
	assert.NotEqual(t, "", fileContents)
}

func TestRenderingEngine_WriteFile(t *testing.T) {
	_, engine := initRepository(t, FileRenderTestDefault, "default")
	fileToRender := &config_generic.FileToRender{
		FileIn: "test_fs/templates/render-context.json",
		FileOut: "render-context.json",
	}
	engine.fsOut = afero.NewMemMapFs()
	_, err := engine.WriteFile(fileToRender)
	assert.NoError(t, err)

	fileExists, errExists := afero.Exists(engine.fsOut, fileToRender.FileOut)
	assert.NoError(t, errExists)
	assert.True(t, fileExists)

}

func TestRenderingEngine_createTemplate(t *testing.T) {
	fs := afero.FromIOFS{FS: testFiles}

	t.Run("create template for existing file", func(t *testing.T) {
		tpl, err := createTemplate(fs, "test_fs/templates/render-context.json")
		assert.NoError(t, err)
		assert.NotNil(t, tpl)
	})

	t.Run("fail if file not exists", func(t *testing.T) {
		tpl, err := createTemplate(fs, "test_fs/templates/missing-file")
		assert.Error(t, err)
		assert.Nil(t, tpl)
	})


}
