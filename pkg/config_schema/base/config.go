package base

import (
	"encoding/base64"
	"fmt"
	"github.com/benammann/git-secrets/pkg/encryption"
)

type Context struct {
	Name string
	SecretResolver encryption.SecretResolver
	EncryptedSecrets map[string]string
	Encryption encryption.Engine
	Secrets []*Secret
}

func (c *Context) EncodeValue(plainValue string) (encodedValue string, err error) {
	encodedString, errEncode := c.Encryption.EncodeValue(plainValue)
	if errEncode != nil {
		return "", errEncode
	}
	return base64.StdEncoding.EncodeToString([]byte(encodedString)), nil
}

func (c *Context) DecodeValue(encodedValue string) (decodedValue string, err error) {
	decodedBase64Bytes, errB64 := base64.StdEncoding.DecodeString(encodedValue)
	if errB64 != nil {
		return "", fmt.Errorf("could not decode base64 value: %s", errB64.Error())
	}
	decodedString, errDecode := c.Encryption.DecodeValue(string(decodedBase64Bytes))
	if errDecode != nil {
		return "", errDecode
	}
	return decodedString, nil
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

type ConfigCliArgs struct {
	OverwriteSecret string
	OverwriteSecretName string
	OverwriteSecretEnv string
}

func (c *Config) GetDefaultContext() *Context {
	return c.GetContext("default")
}

func (c *Config) GetContext(contextName string) *Context {
	for _, context := range c.Contexts {
		if context.Name == contextName {
			return context
		}
	}
	return nil
}

func (c *Config) MergeWithCliArgs(overwrites ConfigCliArgs) {

	var newSecretResolver encryption.SecretResolver
	if overwrites.OverwriteSecret != "" {
		newSecretResolver = encryption.NewPlainSecretResolver(overwrites.OverwriteSecret)
	} else if(overwrites.OverwriteSecretName) != "" {
		newSecretResolver = encryption.NewNameSecretResolver(overwrites.OverwriteSecretName)
	} else if(overwrites.OverwriteSecretEnv) != "" {
		newSecretResolver = encryption.NewEnvSecretResolver(overwrites.OverwriteSecretEnv)
	}

	if newSecretResolver != nil {
		for _, context := range c.Contexts {
			context.SecretResolver = newSecretResolver
		}
	}

}