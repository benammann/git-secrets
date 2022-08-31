package config_generic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseRepository(t *testing.T) {

	t.Run("fail if absolute file name could not be created", func(t *testing.T) {
		// todo: make filepath.abs fail
	})

	t.Run("fail if file cannot be read from fs", func(t *testing.T) {
		_, errParse := createTestRepository("missing-file.json", "default")
		assert.Error(t, errParse)
	})

	t.Run("fail if invalid json is passed", func(t *testing.T) {
		_, errParse := createTestRepository(TestFileInvalidJson, "default")
		assert.Error(t, errParse)
	})

	t.Run("fail if non supported version number is passed", func(t *testing.T) {
		_, errParse := createTestRepository(TestFileInvalidVersion, "default")
		assert.Error(t, errParse)
	})

	t.Run("v1: fail if error while parsing", func(t *testing.T) {
		_, errParse := createTestRepository(TestFileInvalidJsonV1, "default")
		assert.Error(t, errParse)
	})

	t.Run("v1: parse valid config", func(t *testing.T) {
		_, errParse := createTestRepository(TestFileRealWorld, "default")
		assert.NoError(t, errParse)
	})

}
