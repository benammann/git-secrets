package base

import "github.com/benammann/git-secrets/pkg/encryption"

type Context struct {
	Name string
	SecretResolver encryption.SecretResolver
	EncryptedSecrets map[string]string
}

type RenderFilesContext struct {
	ContextName string
	Files []*RenderFilesFile
}

type RenderFilesFile struct {
	FileIn string
	FileOut string
}

type Config struct {
	Contexts []*Context
	RenderFileContexts []*RenderFilesContext
}

func (c *Config) GetDefaultContext() *Context {
	for _, context := range c.Contexts {
		if context.Name == "default" {
			return context
		}
	}
	return nil
}