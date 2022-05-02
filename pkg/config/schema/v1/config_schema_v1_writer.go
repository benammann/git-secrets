package config_schema_v1

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type V1Writer struct {
	schema     Schema
	configPath string
}

func NewV1Writer(schema Schema, configPath string) *V1Writer {
	return &V1Writer{
		schema:     schema,
		configPath: configPath,
	}
}

func (v *V1Writer) AddSecret(contextName string, secretName string, secretEncodedValue string) error {

	if v.schema.Context[contextName] == nil {
		return fmt.Errorf("the context %s does not exist", contextName)
	}

	if v.schema.Context[contextName].Secrets == nil {
		v.schema.Context[contextName].Secrets = make(map[string]string)
	}

	if v.schema.Context[contextName].Secrets[secretName] != "" {
		return fmt.Errorf("the secret %s does already exist. Please overwrite manually", secretName)
	}

	v.schema.Context[contextName].Secrets[secretName] = secretEncodedValue

	return v.WriteConfig()

}

func (v *V1Writer) AddContext(contextName string) error {

	if v.schema.Context[contextName] != nil {
		return fmt.Errorf("the context %s does already exist", contextName)
	}

	v.schema.Context[contextName] = &ContextAwareSecrets{
		Secrets: make(map[string]string),
		Configs: make(map[string]string),
	}

	return v.WriteConfig()

}

func (v *V1Writer) AddFileToRender(contextName string, fileIn string, fileOut string) error {

	absFileIn, _ := filepath.Abs(fileIn)
	absFileOut, _ := filepath.Abs(fileIn)

	if v.schema.RenderFiles == nil {
		v.schema.RenderFiles = make(map[string]*ContextAwareFilesToRender)
	}

	if v.schema.RenderFiles[contextName] == nil {
		v.schema.RenderFiles[contextName] = &ContextAwareFilesToRender{
			Files: []*ContextAwareFileEntry{},
		}
	}

	fileAlreadyAdded := false
	for _, fileToRender := range v.schema.RenderFiles[contextName].Files {

		absFileInCurrent, _ := filepath.Abs(fileToRender.FileIn)
		absFileOutCurrent, _ := filepath.Abs(fileToRender.FileOut)

		if absFileInCurrent == absFileIn && absFileOutCurrent == absFileOut {
			fileAlreadyAdded = true
			break
		}

	}

	if fileAlreadyAdded {
		return fmt.Errorf("the combination %s / %s is already added to context %s", fileIn, fileOut, contextName)
	}

	v.schema.RenderFiles[contextName].Files = append(v.schema.RenderFiles[contextName].Files, &ContextAwareFileEntry{
		FileIn:  fileIn,
		FileOut: fileOut,
	})

	return v.WriteConfig()

}

func (v *V1Writer) WriteConfig() error {

	if errValidate := v.schema.validate(); errValidate != nil {
		return fmt.Errorf("not writing config since it is not valid: %s", errValidate.Error())
	}

	newConfig, _ := json.MarshalIndent(v.schema, "", "  ")

	f, err := os.OpenFile(v.configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return fmt.Errorf("could not open config: %s", err.Error())
	}

	defer f.Close()

	_, errWrite := f.Write(newConfig)
	if errWrite != nil {
		return fmt.Errorf("could not overwrite config: %s", errWrite.Error())
	}

	return nil

}
