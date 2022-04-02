package config_schema_v1

import (
	"encoding/json"
	"fmt"
	config_const "github.com/benammann/git-secrets/pkg/config/const"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	"github.com/benammann/git-secrets/pkg/encryption"
	"github.com/benammann/git-secrets/schema"
	"github.com/xeipuuv/gojsonschema"
	"sort"
)

type Schema struct {
	Version     int                                  `json:"version"`
	Context     Context                              `json:"context"`
	RenderFiles map[string]ContextAwareFilesToRender `json:"renderFiles"`
}

type DecryptSecret struct {
	FromName string `json:"fromName"`
	FromEnv  string `json:"fromEnv"`
}

type ContextAwareSecrets struct {
	DecryptSecret DecryptSecret     `json:"decryptSecret"`
	Secrets       map[string]string `json:"secrets"`
}

type Context map[string]*ContextAwareSecrets

type ContextAwareFileEntry struct {
	FileIn  string `json:"fileIn"`
	FileOut string `json:"fileOut"`
}

type ContextAwareFilesToRender struct {
	Files []ContextAwareFileEntry `json:"files"`
}

type Features struct {
}

var jsonLoader gojsonschema.JSONLoader

func init() {
	jsonLoader = gojsonschema.NewStringLoader(string(schema.GetSchemaContents(schema.V1)))
}

func IsV1(version int) bool {
	return !(version < 1 || version > 1)
}

func (s *Schema) validate() error {

	if !IsV1(s.Version) {
		return fmt.Errorf("not able to process version %d", s.Version)
	}

	// check for default context
	if s.Context["default"] == nil {
		return fmt.Errorf("context.default is required")
	}

	// check for only one or none decryptSecret method
	for contextKey, contextValue := range s.Context {
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

		for secretKey, _ := range contextValue.Secrets {
			if defaultContext.Secrets[secretKey] == "" {
				return fmt.Errorf("secret %s exists in context %s but not in default context", secretKey, contextKey)
			}
		}

	}

	if s.RenderFiles != nil {
		for renderFilesContextKey, _ := range s.RenderFiles {
			if s.Context[renderFilesContextKey] == nil {
				return fmt.Errorf("context %s is defined in features.renderFiles but not in context.%s", renderFilesContextKey, renderFilesContextKey)
			}
		}
	}

	return nil

}

func ParseSchemaV1(jsonInput []byte) (*config_generic.Repository, error) {

	repository := config_generic.NewRepository(1)

	jsonContentLoader := gojsonschema.NewStringLoader(string(jsonInput))
	res, errValidate := gojsonschema.Validate(jsonLoader, jsonContentLoader)
	if errValidate != nil {
		return nil, fmt.Errorf("could not validate json schema: %s", errValidate.Error())
	}

	if res.Valid() == false {
		for _, schemaErr := range res.Errors() {
			fmt.Println(schemaErr.String())
		}
		return nil, fmt.Errorf("invalid json passed")
	}

	var Parsed Schema
	errParse := json.Unmarshal(jsonInput, &Parsed)
	if errParse != nil {
		return nil, fmt.Errorf("could not parse json: %s", errParse.Error())
	}

	if errValidate := Parsed.validate(); errValidate != nil {
		return nil, fmt.Errorf("validation error: %s", errValidate.Error())
	}

	// all resulting contexts
	var contexts []*config_generic.Context

	// locally store the default context
	var defaultContext *config_generic.Context

	//var renderFileContexts []*base.RenderFilesContext

	// first, initialize all contexts
	for contextKey, contextValue := range Parsed.Context {
		localContext := &config_generic.Context{
			Name:             contextKey,
			EncryptedSecrets: contextValue.Secrets,
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
		context.SecretResolver = getSecretResolver(Parsed.Context[context.Name].DecryptSecret, defaultContext)
		context.Encryption = encryption.NewAesEngine(context.SecretResolver)
	}

	if Parsed.RenderFiles != nil {
		for _, context := range contexts {
			if Parsed.RenderFiles[context.Name].Files != nil {
				for _, fileToRender := range Parsed.RenderFiles[context.Name].Files {
					errAddFile := context.AddFileToRender(fileToRender.FileIn, fileToRender.FileOut)
					if errAddFile != nil {
						return nil, fmt.Errorf("could not add file (%s -> %s) to context %s: %s", fileToRender.FileIn, fileToRender.FileOut, context.Name, errAddFile.Error())
					}
				}
			}
		}
	}

	var secrets []*config_generic.Secret

	for _, context := range contexts {
		for secretKey, encryptedSecret := range context.EncryptedSecrets {
			secrets = append(secrets, &config_generic.Secret{
				Name:          secretKey,
				OriginContext: context,
				EncodedValue:  encryptedSecret,
			})
		}
	}

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

	return repository, nil

}

func getSecretResolver(val DecryptSecret, defaultContext *config_generic.Context) encryption.SecretResolver {
	if val.FromEnv != "" {
		return encryption.NewEnvSecretResolver(val.FromEnv)
	}
	if val.FromName != "" {
		return encryption.NewNameSecretResolver(val.FromName)
	}
	return defaultContext.SecretResolver
}
