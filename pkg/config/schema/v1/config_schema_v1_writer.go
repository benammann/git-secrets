package config_schema_v1

import (
	"encoding/json"
	"fmt"
	config_const "github.com/benammann/git-secrets/pkg/config/const"
	"os"
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

func (v *V1Writer) SetSecret(contextName string, secretName string, secretEncodedValue string, force bool) error {

	if v.schema.Context[contextName] == nil {
		return fmt.Errorf("the context %s does not exist", contextName)
	}

	if v.schema.Context[contextName].Secrets == nil {
		v.schema.Context[contextName].Secrets = make(map[string]string)
	}

	if v.schema.Context[contextName].Secrets[secretName] != "" && force == false {
		return fmt.Errorf("the secret %s does already exist. Use --force to overwrite", secretName)
	}

	v.schema.Context[contextName].Secrets[secretName] = secretEncodedValue

	return v.WriteConfig()

}

func (v *V1Writer) SetConfig(contextName string, configName string, configValue string, force bool) error {

	if v.schema.Context[contextName] == nil {
		return fmt.Errorf("the context %s does not exist. Use git-secrets add context <contextName> to add a context", contextName)
	}

	if v.schema.Context[contextName].Configs == nil {
		v.schema.Context[contextName].Configs = make(map[string]string)
	}

	if v.schema.Context[contextName].Configs[configName] != "" && force == false {
		return fmt.Errorf("the config entry %s does already exist. Use --force to overwrite", configName)
	}

	v.schema.Context[contextName].Configs[configName] = configValue

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
		if fileToRender.FileIn == fileIn && fileToRender.FileOut == fileOut {
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

	for contextName, context := range v.schema.Context {
		if context.Secrets == nil && contextName == config_const.DefaultContextName {
			context.Secrets = make(map[string]string)
		}
	}

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
