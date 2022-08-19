package config_generic

import (
	"encoding/json"
	"fmt"
	global_config "github.com/benammann/git-secrets/pkg/config/global"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
)

type VersionFixType struct {
	Version int `json:"version"`
}

func ParseRepositoryFromReadFileFs(fileSystem fs.ReadFileFS, fileName string, globalConfig *global_config.GlobalConfigProvider, overwrittenSecrets map[string]string) (*Repository, error) {

	pathToFile, errAbs := filepath.Abs(fileName)
	if errAbs != nil {
		return nil, fmt.Errorf("could not create absolute file path of config file: %s", errAbs.Error())
	}

	fileContents, fileErr := fileSystem.ReadFile(fileName)
	if fileErr != nil {
		log.Fatalf("could not load test file %s: %s", fileName, fileErr.Error())
	}

	return parseRepositoryFromContents(fileContents, pathToFile, globalConfig, overwrittenSecrets)
}

func ParseRepositoryFromPath(pathToFile string, globalConfig *global_config.GlobalConfigProvider, overwrittenSecrets map[string]string) (*Repository, error) {
	pathToFile, errAbs := filepath.Abs(pathToFile)
	if errAbs != nil {
		return nil, fmt.Errorf("could not create absolute file path of config file: %s", errAbs.Error())
	}

	fileContents, errRead := ioutil.ReadFile(pathToFile)
	if errRead != nil {
		return nil, fmt.Errorf("could not open file: %s", errRead.Error())
	}

	return parseRepositoryFromContents(fileContents, pathToFile, globalConfig, overwrittenSecrets)

}

func parseRepositoryFromContents(fileContents []byte, pathToFile string, globalConfig *global_config.GlobalConfigProvider, overwrittenSecrets map[string]string) (*Repository, error) {
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
