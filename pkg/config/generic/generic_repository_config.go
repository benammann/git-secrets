package config_generic

import (
	"fmt"
	config_const "github.com/benammann/git-secrets/pkg/config/const"
	"sort"
)

type Config struct {

	// Name describes the name of the config entry
	Name string

	// Value hold the encodedValue in base64 of the secret
	Value string

	// OriginContext references the configured context to decode the secret
	OriginContext *Context
}

// AddConfig adds a Config to the repository
// also does some validations
func (c *Repository) AddConfig(Config *Config) error {

	// if not default Config we need to check if the given Config is also configured in the default context
	// because we are not allowed to define variables only in a child context
	if Config.OriginContext.Name != config_const.DefaultContextName {

		// get all the default Configs
		defaultConfigs := c.GetConfigsByContext(config_const.DefaultContextName)
		defaultConfigFound := false

		// check if it defined
		for _, defaultConfig := range defaultConfigs {
			if defaultConfig.Name == Config.Name {
				defaultConfigFound = true
				break
			}
		}

		// return error if not defined
		if defaultConfigFound == false {
			return fmt.Errorf("config %s defined in context %s is not defined in the default context", Config.Name, Config.OriginContext.Name)
		}

	}

	// append the Config to the repository
	c.configs = append(c.configs, Config)

	// sort the Configs alphabetically
	sort.SliceStable(c.configs, func(i, j int) bool {
		return c.configs[i].Name < c.configs[j].Name
	})

	return nil
}

// GetConfigsByContext returns all the Configs related to the current context
func (c *Repository) GetConfigsByContext(contextName string) (res []*Config) {
	for _, config := range c.configs {
		if config.OriginContext.Name == contextName {
			res = append(res, config)
		}
	}
	return res
}

// GetCurrentConfigs merges the default Configs with the context Configs
// the default Configs are overwritten by the context Configs
func (c *Repository) GetCurrentConfigs() (res []*Config) {

	// get all default Configs
	defaultConfigs := c.GetConfigsByContext(config_const.DefaultContextName)

	// if not default, merge the Configs with the default Configs
	if !c.IsDefault() {

		// result is context Configs
		contextConfigs := c.GetConfigsByContext(c.context.Name)

		// add the default Config if it is missing in the context Configs
		for _, defaultConfig := range defaultConfigs {

			found := false
			for _, contextConfig := range contextConfigs {
				if defaultConfig.Name == contextConfig.Name {
					found = true
					break
				}
			}

			if found == false {
				contextConfigs = append(contextConfigs, defaultConfig)
			}

		}

		res = contextConfigs

	} else {
		// result is just the default Configs
		res = defaultConfigs
	}

	return res

}

// GetCurrentConfig takes the merged Configs from GetCurrentConfigs and returns the needed one
func (c *Repository) GetCurrentConfig(ConfigName string) *Config {
	for _, Config := range c.GetCurrentConfigs() {
		if Config.Name == ConfigName {
			return Config
		}
	}
	return nil
}

func (c *Repository) GetConfigMap() ConfigMap {
	mapOut := make(ConfigMap)
	configValues := c.GetCurrentConfigs()
	for _, configItem := range configValues {
		mapOut[configItem.Name] = configItem.Value
	}
	return mapOut
}
