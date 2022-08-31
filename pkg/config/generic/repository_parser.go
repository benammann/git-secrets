package config_generic

import (
	"encoding/json"
	"fmt"
	global_config "github.com/benammann/git-secrets/pkg/config/global"
	"github.com/spf13/afero"
	"path/filepath"
)

type VersionFixType struct {
	Version int `json:"version"`
}

func ParseRepository(fileSystem afero.Fs, fileName string, globalConfig *global_config.GlobalConfigProvider, overwrittenSecrets map[string]string) (*Repository, error) {

	pathToFile, _ := filepath.Abs(fileName)

	fileContents, fileErr := afero.ReadFile(fileSystem, fileName)
	if fileErr != nil {
		return nil, fmt.Errorf("could not load test file %s: %s", fileName, fileErr.Error())
	}

	var VersionBase VersionFixType
	errParse := json.Unmarshal(fileContents, &VersionBase)
	if errParse != nil {
		return nil, fmt.Errorf("could not parse json: %s", errParse.Error())
	}

	version := VersionBase.Version
	if IsSchemaV1(version) {
		configOut, errConfigOut := ParseSchemaV1(fileContents, pathToFile, globalConfig, overwrittenSecrets)
		if errConfigOut != nil {
			return nil, fmt.Errorf("could not parse as v1: %s", errConfigOut.Error())
		}
		return configOut, nil
	} else {
		return nil, fmt.Errorf("unsupported version: %d", version)
	}

}
