package config_generic

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContext_DecodeValue(t *testing.T) {

	repo, errParse := createTestRepository(TestFileBlankDefault, "default")
	assert.NotNil(t, repo)
	assert.NoError(t, errParse)

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

	repo, errParse := createTestRepository(TestFileBlankDefault, "default")
	assert.NotNil(t, repo)
	assert.NoError(t, errParse)

	ctx := repo.GetCurrent()
	encodedValue, errEncode := ctx.EncodeValue("Hello World")
	assert.NoError(t, errEncode)

	_, errB64 := base64.StdEncoding.DecodeString(encodedValue)
	assert.NoError(t, errB64)

}

func TestRepository_AddContext(t *testing.T) {

}

func TestRepository_GetContext(t *testing.T) {

}

func TestRepository_GetContexts(t *testing.T) {

}

func TestRepository_GetCurrent(t *testing.T) {

}

func TestRepository_GetDefault(t *testing.T) {

}

func TestRepository_SetSelectedContext(t *testing.T) {

}
