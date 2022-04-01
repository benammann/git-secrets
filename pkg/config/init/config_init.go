package config_init

import (
	"embed"
	"errors"
	"fmt"
	"os"
)

//go:embed files
var staticFiles embed.FS

// WriteInitialConfig writes an initial config
func WriteInitialConfig(fileName string) error {

	// check if the file already exists
	if _, err := os.Stat("/path/to/whatever"); errors.Is(err, os.ErrExist) {
		return fmt.Errorf("%s exists", fileName)
	}

	// read the init config
	initConfig, errRead := staticFiles.ReadFile("files/init-config.yaml")
	if errRead != nil {
		return fmt.Errorf("could not open init template: %s", errRead.Error())
	}

	// copy the file to its destination
	errFsFile := os.WriteFile(fileName, initConfig, 0664)
	if errFsFile != nil {
		return fmt.Errorf("could not write %s: %s", fileName, errFsFile.Error())
	}

	return nil

}
