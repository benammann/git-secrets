package encryption

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func newTestAesEngine(t *testing.T) *AesEngine {
	assert.NoError(t, os.Setenv("SR_ENV", "aju1ZieThohngii4eem4saeCh2fieral"))
	sr := NewEnvSecretResolver("SR_ENV")
	return NewAesEngine(sr)
}

func TestAesEngine_DecodeValue(t *testing.T) {
	engine := newTestAesEngine(t)
	t.Run("fail if unable to decode string", func(t *testing.T) {
		_, errDecode := engine.DecodeValue("abcdefg")
		assert.Error(t, errDecode)
	})
	t.Run("decode encrypted values", func(t *testing.T) {
		str := "hello world"
		encodedValue, errEncode := engine.EncodeValue(str)
		assert.NoError(t, errEncode)
		decodedValue, errDecode := engine.DecodeValue(encodedValue)
		assert.NoError(t, errDecode)
		assert.Equal(t, str, decodedValue)
	})
}

func TestAesEngine_EncodeValue(t *testing.T) {
	engine := newTestAesEngine(t)
	t.Run("encode values", func(t *testing.T) {
		str := "hello world"
		encodedValue, errEncode := engine.EncodeValue(str)
		assert.NoError(t, errEncode)
		decodedValue, errDecode := engine.DecodeValue(encodedValue)
		assert.NoError(t, errDecode)
		assert.Equal(t, str, decodedValue)
	})
}

func TestAesEngine_newGcm(t *testing.T) {
	engine := newTestAesEngine(t)
	_, _, errGcm := engine.newGcm()
	assert.NoError(t, errGcm)
}

func TestNewAesEngine(t *testing.T) {
	newTestAesEngine(t)
}
