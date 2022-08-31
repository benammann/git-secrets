package main

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseEnv(t *testing.T) {

	fs := afero.NewMemMapFs()
	assert.NoError(t, afero.WriteFile(fs, ".env", []byte(`
ENV_A=ENV_A_VALUE
# a comment
ENV_B=ENV_B_VALUE
# ENV_C=ENV_C_VALUE
`), 0664))

	assert.NoError(t, ParseEnv(fs, ".env"))
	assert.Equal(t, "ENV_A_VALUE", os.Getenv("ENV_A"))
	assert.Equal(t, "ENV_B_VALUE", os.Getenv("ENV_B"))
	assert.Equal(t, "", os.Getenv("ENV_C"))

	assert.Error(t, ParseEnv(fs, ".env-missing"))

}
