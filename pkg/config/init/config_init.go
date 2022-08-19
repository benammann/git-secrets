package config_init

import (
	"embed"
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"os"
	"strings"
)

//go:embed files
var staticFiles embed.FS

// WriteInitialConfig writes an initial config
func WriteInitialConfig(fileSystem afero.Fs, fileName string, secretName string) error {

	// check if the file already exists
	if _, err := fileSystem.Stat(fileName); errors.Is(err, os.ErrExist) {
		return fmt.Errorf("%s exists", fileName)
	}

	// read the init config
	initConfig, errRead := staticFiles.ReadFile("files/init-config.json")
	if errRead != nil {
		return fmt.Errorf("could not open init template: %s", errRead.Error())
	}

	finalInitConfig := strings.ReplaceAll(string(initConfig), "{{secretName}}", secretName)

	// copy the file to its destination
	errFsFile := afero.WriteFile(fileSystem, fileName, []byte(finalInitConfig), 0664)
	if errFsFile != nil {
		return fmt.Errorf("could not write %s: %s", fileName, errFsFile.Error())
	}

	return nil

}
