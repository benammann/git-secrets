package config_generic

import (
	"encoding/base64"
	"fmt"
	config_const "github.com/benammann/git-secrets/pkg/config/const"
	"github.com/benammann/git-secrets/pkg/encryption"
)

type Context struct {
	Name             string
	SecretResolver   encryption.SecretResolver
	Encryption       encryption.Engine
	EncryptedSecrets map[string]string
	Configs          map[string]string
	FilesToRender    []*FileToRender
}

type FileToRender struct {
	FileIn  string
	FileOut string
}

// AddContext adds a context and does some validations
func (c *Repository) AddContext(context *Context) error {
	defaultContext := c.GetDefault()
	if defaultContext == nil && context.Name != config_const.DefaultContextName {
		return fmt.Errorf("the default context must be added first")
	}
	if defaultContext != nil && context.Name == config_const.DefaultContextName {
		return fmt.Errorf("the default context is already defined")
	}
	c.contexts = append(c.contexts, context)
	return nil
}

// GetContext returns the context by name
func (c *Repository) GetContext(contextName string) *Context {
	for _, context := range c.contexts {
		if context.Name == contextName {
			return context
		}
	}
	return nil
}

func (c *Repository) GetContexts() []*Context {
	return c.contexts
}

// GetDefault returns the default context
func (c *Repository) GetDefault() *Context {
	return c.GetContext(config_const.DefaultContextName)
}

// SetSelectedContext sets the current selected context
func (c *Repository) SetSelectedContext(contextName string) (*Context, error) {
	desiredContext := c.GetContext(contextName)
	if desiredContext == nil {
		return nil, fmt.Errorf("the context %s does not exists", contextName)
	}
	c.context = desiredContext
	return desiredContext, nil
}

// GetCurrent returns the current context
func (c *Repository) GetCurrent() *Context {
	return c.context
}

// EncodeValue encodes the given value and returns it as a base64 string
func (c *Context) EncodeValue(plainValue string) (encodedValue string, err error) {
	encodedString, errEncode := c.Encryption.EncodeValue(plainValue)
	if errEncode != nil {
		return "", errEncode
	}
	return base64.StdEncoding.EncodeToString([]byte(encodedString)), nil
}

// DecodeValue takes the value and decodes it
// encodeValue must be base64
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

// AddFileToRender adds a file to render which is later used by the rendering engine
func (c *Context) AddFileToRender(fileIn string, fileOut string) error {

	// check if output file is double defined
	for _, fileToRender := range c.FilesToRender {
		if fileToRender.FileOut == fileOut {
			return fmt.Errorf("output file %s is already defined on context %s", fileOut, c.Name)
		}
	}

	c.FilesToRender = append(c.FilesToRender, &FileToRender{
		FileIn:  fileIn,
		FileOut: fileOut,
	})

	return nil
}
