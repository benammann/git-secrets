package config_parser

import (
	"fmt"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	config_schema_v1 "github.com/benammann/git-secrets/pkg/config/schema/v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type VersionFixType struct {
	Version int `yaml:"version"`
}

func ParseRepository(pathToFile string) (*config_generic.Repository, error) {

	fileContents, errRead := ioutil.ReadFile(pathToFile)
	if errRead != nil {
		return nil, fmt.Errorf("could not open file: %s", errRead.Error())
	}

	var VersionBase VersionFixType
	errParse := yaml.Unmarshal(fileContents, &VersionBase)
	if errParse != nil {
		return nil, fmt.Errorf("could not parse yaml: %s", errParse.Error())
	}

	version := VersionBase.Version
	if config_schema_v1.IsV1(version) {
		configOut, errConfigOut := config_schema_v1.ParseSchemaV1(fileContents)
		if errConfigOut != nil {
			return nil, fmt.Errorf("could not parse as v1: %s", errConfigOut.Error())
		}
		return configOut, nil
	} else {
		return nil, fmt.Errorf("unsupported version: %d", version)
	}

}
