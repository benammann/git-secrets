package config_init

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

const OutputFile = ".git-secrets.json"
const DecryptSecretName = "myLocalSecret"

func TestWriteInitialConfig(t *testing.T) {

	var fs = afero.NewMemMapFs()

	t.Run("write config file with the correct decrypt secret", func(t *testing.T) {

		assert.NoError(t, WriteInitialConfig(fs, OutputFile, DecryptSecretName))

		_, fileErr := fs.Stat(OutputFile)
		assert.NoError(t, fileErr)

		file, errOpen := fs.Open(OutputFile)
		defer file.Close()
		assert.NoError(t, errOpen)

		fileBytes, errRead := ioutil.ReadAll(file)
		assert.NoError(t, errRead)

		assert.True(t, strings.Contains(string(fileBytes), fmt.Sprintf("\"fromName\": \"%s\"", DecryptSecretName)))

	})

}
