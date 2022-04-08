package config_generic

import (
	config_const "github.com/benammann/git-secrets/pkg/config/const"
)

type SecretsMap map[string]string
type ConfigMap map[string]string

// NewRepository creates a new generic repository
func NewRepository(configVersion int, configFileUsed string) *Repository {
	return &Repository{
		configVersion:  configVersion,
		configFileUsed: configFileUsed,
	}
}

type Repository struct {

	// configVersion hold the config version this repository is built from
	configVersion int

	// configFileUsed holds the abs path of the used config file
	configFileUsed string

	// context holds the current resolved context
	context *Context

	// contexts holds all the available contexts
	contexts []*Context

	// secrets holds all secrets of all contexts
	secrets []*Secret

	// configs holds all configs of all contexts
	configs []*Config
}

// GetConfigVersion returns the config version this repository is built from
func (c *Repository) GetConfigVersion() int {
	return c.configVersion
}

// IsDefault returns if the default context is used
func (c *Repository) IsDefault() bool {
	return c.context.Name == config_const.DefaultContextName
}
