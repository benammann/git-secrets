package config_generic

import "github.com/benammann/git-secrets/pkg/encryption"

type ConfigCliArgs struct {
	OverwriteSecret string
	OverwriteSecretName string
	OverwriteSecretEnv string
}

// MergeWithCliArgs overwrites the repository config with cli arguments
func (c *Repository) MergeWithCliArgs(overwrites ConfigCliArgs) {

	var newSecretResolver encryption.SecretResolver
	if overwrites.OverwriteSecret != "" {
		newSecretResolver = encryption.NewPlainSecretResolver(overwrites.OverwriteSecret)
	} else if(overwrites.OverwriteSecretName) != "" {
		newSecretResolver = encryption.NewNameSecretResolver(overwrites.OverwriteSecretName)
	} else if(overwrites.OverwriteSecretEnv) != "" {
		newSecretResolver = encryption.NewEnvSecretResolver(overwrites.OverwriteSecretEnv)
	}

	if newSecretResolver != nil {
		for _, context := range c.GetContexts() {
			context.SecretResolver = newSecretResolver
		}
	}

}
