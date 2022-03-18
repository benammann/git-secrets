package v1

import (
	"fmt"
	"github.com/benammann/git-secrets/pkg/config_schema/base"
	"github.com/benammann/git-secrets/pkg/encryption"
	"gopkg.in/yaml.v2"
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
	DecryptSecret DecryptSecret `yaml:"decryptSecret"`
	Secrets       map[string]string       `yaml:"secrets"`
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

func ParseSchemaV1(input []byte) (*base.Config, error) {

	var Parsed Schema
	errParse := yaml.Unmarshal(input, &Parsed)
	if errParse != nil {
		return nil, fmt.Errorf("could not parse yaml: %s", errParse.Error())
	}

	if errValidate := Parsed.validate(); errValidate != nil {
		return nil, fmt.Errorf("validation error: %s", errValidate.Error())
	}

	var contexts []*base.Context
	var renderFileContexts []*base.RenderFilesContext
	var defaultContext *base.Context

	for contextKey, contextValue := range Parsed.Context {

		localContext := &base.Context{
			Name: contextKey,
			EncryptedSecrets: contextValue.Secrets,
		}

		if localContext.Name == "default" {
			defaultContext = localContext
		}

		contexts = append(contexts, localContext)
	}

	for _, context := range contexts {
		context.SecretResolver = getSecretResolver(Parsed.Context[context.Name].DecryptSecret, defaultContext)
	}

	for contextKey, contextValue := range Parsed.Features.RenderFiles {

		var files []*base.RenderFilesFile
		for _, file := range contextValue.Files {
			files = append(files, &base.RenderFilesFile{
				FileIn: file.FileIn,
				FileOut: file.FileOut,
			})
		}

		renderFileContexts = append(renderFileContexts, &base.RenderFilesContext{
			ContextName: contextKey,
			Files: files,
		})
	}

	out := &base.Config{
		Contexts: contexts,
		RenderFileContexts: renderFileContexts,
	}

	return out, nil

}

func getSecretResolver(val DecryptSecret, defaultContext *base.Context) encryption.SecretResolver {
	if val.FromEnv != "" {
		return encryption.NewEnvSecretResolver(val.FromEnv)
	}
	if val.FromName != "" {
		return encryption.NewNameSecretResolver(val.FromName)
	}
	return defaultContext.SecretResolver
}
