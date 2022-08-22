package config_generic

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContext_DecodeValue(t *testing.T) {

	repo := initRepository(t, TestFileBlankDefault, "default")

	ctx := repo.GetCurrent()

	t.Run("it should decode valid base64", func(t *testing.T) {
		decodedValue, errDecode := ctx.DecodeValue("l5sqnu8UkO+PdW2fZo7IMhfHng7lf6XNXEfRhQ/fvboP1HqcRFcu")
		assert.NoError(t, errDecode)
		assert.Equal(t, "Hello World", decodedValue)
	})

	t.Run("it should not decode invalid base64", func(t *testing.T) {
		decodedValue, errDecode := ctx.DecodeValue("abc")
		assert.Error(t, errDecode)
		assert.Equal(t, "", decodedValue)
	})

	t.Run("it should not decode invalid values", func(t *testing.T) {
		decodedValue, errDecode := ctx.DecodeValue("YWJjCg==")
		assert.Error(t, errDecode)
		assert.Equal(t, "", decodedValue)
	})

}

func TestContext_EncodeValue(t *testing.T) {

	repo := initRepository(t, TestFileBlankDefault, "default")

	ctx := repo.GetCurrent()
	encodedValue, errEncode := ctx.EncodeValue("Hello World")
	assert.NoError(t, errEncode)

	_, errB64 := base64.StdEncoding.DecodeString(encodedValue)
	assert.NoError(t, errB64)

	failingRepo := initRepository(t, TestFileMissingEncryptionSecret, "default")
	failingCtx := failingRepo.GetCurrent()
	assert.NotNil(t, failingCtx)
	failedValue, expectedErr := failingCtx.EncodeValue("Hello World")
	assert.Error(t, expectedErr)
	assert.Equal(t, "", failedValue)

}

func TestRepository_AddContext(t *testing.T) {

	repo := initRepository(t, TestFileBlankDefault, "default")
	repo.contexts = []*Context{}

	defaultCtx := &Context{
		Name: "default",
	}

	prodCtx := &Context{
		Name: "prod",
	}

	assert.Error(t, repo.AddContext(prodCtx))
	assert.NoError(t, repo.AddContext(defaultCtx))
	assert.Error(t, repo.AddContext(defaultCtx))
	assert.NoError(t, repo.AddContext(prodCtx))

	assert.Equal(t, defaultCtx, repo.GetContext("default"))
	assert.Equal(t, prodCtx, repo.GetContext("prod"))

}

func TestRepository_GetContext(t *testing.T) {

	repo := initRepository(t, TestFileBlankDefault, "default")

	testCtx := &Context{
		Name: "test",
	}

	assert.NoError(t, repo.AddContext(testCtx))

	assert.NotNil(t, repo.GetContext("default"))
	assert.Equal(t, testCtx, repo.GetContext("test"))
	assert.Nil(t, repo.GetContext("missingContext"))

}

func TestRepository_GetContexts(t *testing.T) {
	repo := initRepository(t, TestFileBlankTwoContexts, "default")
	assert.Equal(t, repo.contexts, repo.GetContexts())
	assert.Len(t, repo.GetContexts(), 2)
}

func TestRepository_GetCurrent(t *testing.T) {
	repo := initRepository(t, TestFileBlankTwoContexts, "default")
	defaultCtx := repo.GetContext("default")
	prodCtx := repo.GetContext("prod")
	assert.NotNil(t, defaultCtx)
	assert.NotNil(t, prodCtx)

	assert.Equal(t, defaultCtx, repo.GetCurrent())
	_, errSetCtx := repo.SetSelectedContext("prod")
	assert.NoError(t, errSetCtx)
	assert.Equal(t, prodCtx, repo.GetCurrent())
}

func TestRepository_GetDefault(t *testing.T) {
	repo := initRepository(t, TestFileBlankTwoContexts, "default")
	defaultCtx := repo.GetContext("default")
	assert.NotNil(t, defaultCtx)
	assert.Equal(t, defaultCtx, repo.GetDefault())
}

func TestRepository_SetSelectedContext(t *testing.T) {

	repo := initRepository(t, TestFileBlankTwoContexts, "default")

	defaultCtx := repo.GetContext("default")
	prodCtx := repo.GetContext("prod")
	assert.NotNil(t, defaultCtx)
	assert.NotNil(t, prodCtx)

	ctxOut, errSetCtx := repo.SetSelectedContext("prod")
	assert.NoError(t, errSetCtx)
	assert.Equal(t, ctxOut, prodCtx)

	emptyOut, errSetMissing := repo.SetSelectedContext("missingContext")
	assert.Error(t, errSetMissing)
	assert.Nil(t, emptyOut)

}
