package config

import (
	"fmt"
	"github.com/benammann/git-secrets/pkg/config_schema/base"
	v1 "github.com/benammann/git-secrets/pkg/config_schema/v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type VersionFixType struct {
	Version int `yaml:"version"`
}

func ParseConfig(pathToFile string) (*base.Config, error) {

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
	if v1.IsV1(version) {
		configOut, errConfigOut := v1.ParseSchemaV1(fileContents)
		if errConfigOut != nil {
			return nil, fmt.Errorf("could not parse as v1: %s", errConfigOut.Error())
		}
		return configOut, nil
	} else {
		return nil, fmt.Errorf("unsupported version: %d", version)
	}

}
