package config_generic

import (
	"encoding/json"
	"fmt"
	config_const "github.com/benammann/git-secrets/pkg/config/const"
	global_config "github.com/benammann/git-secrets/pkg/config/global"
	"github.com/benammann/git-secrets/pkg/encryption"
	"github.com/benammann/git-secrets/schema"
	"github.com/spf13/afero"
	"github.com/xeipuuv/gojsonschema"
	"path/filepath"
	"sort"
)

type V1Schema struct {
	Schema      string                     `json:"$schema,omitempty"`
	Version     int                        `json:"version"`
	Context     V1Context                  `json:"context"`
	RenderFiles map[string]*V1RenderTarget `json:"renderFiles,omitempty"`
}

type V1DecryptSecret struct {
	FromName string `json:"fromName,omitempty"`
	FromEnv  string `json:"fromEnv,omitempty"`
}

type V1ContextAwareSecrets struct {
	DecryptSecret *V1DecryptSecret  `json:"decryptSecret,omitempty"`
	Secrets       map[string]string `json:"secrets,omitempty"`
	Configs       map[string]string `json:"configs,omitempty"`
}

type V1Context map[string]*V1ContextAwareSecrets

type V1RenderTargetFileEntry struct {
	FileIn  string `json:"fileIn"`
	FileOut string `json:"fileOut"`
}

type V1RenderTarget struct {
	Files []*V1RenderTargetFileEntry `json:"files"`
}

var jsonLoaderV1 gojsonschema.JSONLoader

func init() {
	jsonLoaderV1 = gojsonschema.NewStringLoader(string(schema.GetSchemaContents(schema.V1)))
}

func IsSchemaV1(version int) bool {
	return !(version < 1 || version > 1)
}

func (s *V1Schema) validateSchemaV1() error {

	if !IsSchemaV1(s.Version) {
		return fmt.Errorf("not able to process version %d", s.Version)
	}

	// check for default context
	if s.Context["default"] == nil {
		return fmt.Errorf("context.default is required")
	}

	// check for only one or none decryptSecret method
	for contextKey, contextValue := range s.Context {
		if contextValue.DecryptSecret == nil {
			continue
		}
		if contextValue.DecryptSecret.FromEnv != "" && contextValue.DecryptSecret.FromName != "" {
			return fmt.Errorf("context: %s: you can only use either one decryptSecret method (FromEnv or FromName)", contextKey)
		}
		if contextValue.DecryptSecret.FromEnv == "" && contextValue.DecryptSecret.FromName == "" && contextKey == "default" {
			return fmt.Errorf("context: %s: you must specify at least one decryption method", contextKey)
		}
	}

	defaultContext := s.Context["default"]

	// check if secret keys exists in default context
	for contextKey, contextValue := range s.Context {

		// skip default context
		if contextKey == "default" {
			continue
		}

		if contextValue.Secrets != nil {
			for secretKey, _ := range contextValue.Secrets {
				if defaultContext.Secrets[secretKey] == "" {
					return fmt.Errorf("secret %s exists in context %s but not in default context", secretKey, contextKey)
				}
			}
		}

		if contextValue.Configs != nil {
			for configKey, _ := range contextValue.Configs {
				if defaultContext.Configs[configKey] == "" {
					return fmt.Errorf("config entry %s exists in context %s but not in default context", configKey, contextKey)
				}
			}
		}

	}

	return nil

}

func ParseSchemaV1(jsonInput []byte, configFileUsed string, globalConfig *global_config.GlobalConfigProvider, overwrittenSecrets map[string]string) (*Repository, error) {

	jsonContentLoader := gojsonschema.NewStringLoader(string(jsonInput))
	res, errValidate := gojsonschema.Validate(jsonLoaderV1, jsonContentLoader)
	if errValidate != nil {
		return nil, fmt.Errorf("could not validateSchemaV1 json schema: %s", errValidate.Error())
	}

	if res.Valid() == false {
		for _, schemaErr := range res.Errors() {
			fmt.Println(schemaErr.String())
		}
		return nil, fmt.Errorf("invalid json passed")
	}

	var Parsed V1Schema
	errParse := json.Unmarshal(jsonInput, &Parsed)
	if errParse != nil {
		return nil, fmt.Errorf("could not parse json: %s", errParse.Error())
	}

	if errValidate := Parsed.validateSchemaV1(); errValidate != nil {
		return nil, fmt.Errorf("validation error: %s", errValidate.Error())
	}

	// all resulting contexts
	var contexts []*Context

	// locally store the default context
	var defaultContext *Context

	// all render targets to add
	var renderTargets []*RenderTarget

	// first, initialize all contexts
	for contextKey, contextValue := range Parsed.Context {
		localContext := &Context{
			Name:             contextKey,
			EncryptedSecrets: contextValue.Secrets,
			Configs:          contextValue.Configs,
		}
		// reference the default context
		if localContext.Name == config_const.DefaultContextName {
			defaultContext = localContext
		}
		contexts = append(contexts, localContext)
	}

	// important, always default first since some logics depend on a fully defined default context
	sort.SliceStable(contexts, func(i, j int) bool {
		return contexts[i].Name == config_const.DefaultContextName
	})

	for _, context := range contexts {
		context.SecretResolver = getSecretResolverV1(Parsed.Context[context.Name].DecryptSecret, defaultContext, globalConfig, overwrittenSecrets)
		context.Encryption = encryption.NewAesEngine(context.SecretResolver)
	}

	if Parsed.RenderFiles != nil {
		for targetName, renderTarget := range Parsed.RenderFiles {
			if renderTarget.Files != nil {
				finalRenderTarget := NewRenderTarget(targetName)
				for _, fileToRender := range renderTarget.Files {
					configDir := filepath.Dir(configFileUsed)
					fileIn := filepath.Join(configDir, fileToRender.FileIn)
					fileOut := filepath.Join(configDir, fileToRender.FileOut)
					errAddFile := finalRenderTarget.AddFileToRender(fileIn, fileOut)
					if errAddFile != nil {
						return nil, fmt.Errorf("could not add file (%s -> %s) to target %s: %s", fileToRender.FileIn, fileToRender.FileOut, finalRenderTarget.Name, errAddFile.Error())
					}
				}
				renderTargets = append(renderTargets, finalRenderTarget)
			}
		}
	}

	var secrets []*Secret

	for _, context := range contexts {
		for secretKey, encryptedSecret := range context.EncryptedSecrets {
			secrets = append(secrets, &Secret{
				Name:          secretKey,
				OriginContext: context,
				EncodedValue:  encryptedSecret,
			})
		}
	}

	var configs []*Config
	for _, context := range contexts {
		for configKey, configValue := range context.Configs {
			configs = append(configs, &Config{
				Name:          configKey,
				Value:         configValue,
				OriginContext: context,
			})
		}
	}

	configWriter := NewV1Writer(afero.NewOsFs(), Parsed, configFileUsed)
	repository := NewRepository(1, configFileUsed, configWriter)

	for _, resultingContext := range contexts {
		errAddContext := repository.AddContext(resultingContext)
		if errAddContext != nil {
			return nil, fmt.Errorf("could not add context to repository: %s", errAddContext.Error())
		}
	}

	for _, secretOut := range secrets {
		errAddSecret := repository.AddSecret(secretOut)
		if errAddSecret != nil {
			return nil, fmt.Errorf("could not add secret to repository: %s", errAddSecret.Error())
		}
	}

	for _, configOut := range configs {
		errAddConfig := repository.AddConfig(configOut)
		if errAddConfig != nil {
			return nil, fmt.Errorf("could not add config to repository: %s", errAddConfig.Error())
		}
	}

	for _, renderTargetOut := range renderTargets {
		if errAddTarget := repository.AddRenderTarget(renderTargetOut); errAddTarget != nil {
			return nil, fmt.Errorf("could not add render target: %s", errAddTarget.Error())
		}
	}

	return repository, nil

}

func getSecretResolverV1(val *V1DecryptSecret, defaultContext *Context, globalConfig *global_config.GlobalConfigProvider, overwrittenSecrets map[string]string) encryption.SecretResolver {
	if val == nil {
		return defaultContext.SecretResolver
	}
	if val.FromEnv != "" {
		return encryption.NewEnvSecretResolver(val.FromEnv)
	}
	if val.FromName != "" {
		return encryption.NewMergedSecretResolver(val.FromName, globalConfig, overwrittenSecrets)
	}
	return defaultContext.SecretResolver
}
