package config_schema_v1

import (
	"fmt"
	config_const "github.com/benammann/git-secrets/pkg/config/const"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	"github.com/benammann/git-secrets/pkg/encryption"
	"gopkg.in/yaml.v2"
	"sort"
)

type Schema struct {
	Version  int      `yaml:"version"`
	Context  Context  `yaml:"context"`
	Features Features `yaml:"features"`
}

type DecryptSecret struct {
	FromName string `yaml:"fromName"`
	FromEnv  string `yaml:"fromEnv"`
}

type ContextAwareSecrets struct {
	DecryptSecret DecryptSecret     `yaml:"decryptSecret"`
	Secrets       map[string]string `yaml:"secrets"`
}

type Context map[string]*ContextAwareSecrets

type ContextAwareFileEntry struct {
	FileIn  string `yaml:"fileIn"`
	FileOut string `yaml:"fileOut"`
}

type ContextAwareFilesToRender struct {
	Files []ContextAwareFileEntry `yaml:"files"`
}

type Features struct {
	RenderFiles map[string]ContextAwareFilesToRender `yaml:"renderFiles"`
	Export      map[string]map[string]string         `yaml:"export"`
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

	for renderFilesContextKey, _ := range s.Features.RenderFiles {
		if s.Context[renderFilesContextKey] == nil {
			return fmt.Errorf("context %s is defined in features.renderFiles but not in context.%s", renderFilesContextKey, renderFilesContextKey)
		}
	}

	return nil

}

func ParseSchemaV1(input []byte) (*config_generic.Repository, error) {

	repository := config_generic.NewRepository(1)

	var Parsed Schema
	errParse := yaml.Unmarshal(input, &Parsed)
	if errParse != nil {
		return nil, fmt.Errorf("could not parse yaml: %s", errParse.Error())
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

	for _, context := range contexts {
		context.SecretResolver = getSecretResolver(Parsed.Context[context.Name].DecryptSecret, defaultContext)
		context.Encryption = encryption.NewAesEngine(context.SecretResolver)
	}

	// for contextKey, contextValue := range Parsed.Features.RenderFiles {

	//var files []*base.RenderFilesFile
	//for _, file := range contextValue.Files {
	//	files = append(files, &base.RenderFilesFile{
	//		FileIn: file.FileIn,
	//		FileOut: file.FileOut,
	//	})
	//}
	//
	//renderFileContexts = append(renderFileContexts, &base.RenderFilesContext{
	//	ContextName: contextKey,
	//	Files: files,
	//})
	// }

	sort.SliceStable(contexts, func(i, j int) bool {
		return contexts[i].Name == "default"
	})

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
