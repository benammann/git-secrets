package config_parser

import (
	"encoding/json"
	"fmt"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	config_schema_v1 "github.com/benammann/git-secrets/pkg/config/schema/v1"
	"io/ioutil"
	"path/filepath"
	yaml2 "sigs.k8s.io/yaml"
)

type VersionFixType struct {
	Version int `json:"version"`
}

func ParseRepository(pathToFile string) (*config_generic.Repository, error) {
	pathToFile, errAbs := filepath.Abs(pathToFile)
	if errAbs != nil {
		return nil, fmt.Errorf("could not create absolute file path of config file: %s", errAbs.Error())
	}

	fileContents, errRead := ioutil.ReadFile(pathToFile)
	if errRead != nil {
		return nil, fmt.Errorf("could not open file: %s", errRead.Error())
	}

	fileExtension := filepath.Ext(pathToFile)

	if fileExtension == ".yaml" || fileExtension == ".yml" {
		jsonBytes, errConvert := yaml2.YAMLToJSON(fileContents)
		if errConvert != nil {
			return nil, fmt.Errorf("could not convert yaml to json")
		}
		fileContents = jsonBytes
	} else if fileExtension == ".json" {
		// file contents is already json
	} else {
		return nil, fmt.Errorf("unsupported file extension: %s. Supported: yaml, json", fileExtension)
	}

	var VersionBase VersionFixType
	errParse := json.Unmarshal(fileContents, &VersionBase)
	if errParse != nil {
		return nil, fmt.Errorf("could not parse json: %s", errParse.Error())
	}

	version := VersionBase.Version
	if config_schema_v1.IsV1(version) {
		configOut, errConfigOut := config_schema_v1.ParseSchemaV1(fileContents, pathToFile)
		if errConfigOut != nil {
			return nil, fmt.Errorf("could not parse as v1: %s", errConfigOut.Error())
		}
		return configOut, nil
	} else {
		return nil, fmt.Errorf("unsupported version: %d", version)
	}

}
