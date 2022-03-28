package config_generic

import (
	"fmt"
	config_const "github.com/benammann/git-secrets/pkg/config/const"
	"sort"
)

type Secret struct {

	// Name describes the name of the secret
	Name string

	// EncodedValue hold the encodedValue in base64 of the secret
	EncodedValue string

	// OriginContext references the configured context to decode the secret
	OriginContext *Context
}

// AddSecret adds a secret to the repository
// also does some validations
func (c *Repository) AddSecret(secret *Secret) error {

	// if not default secret we need to check if the given secret is also configured in the default context
	// because we are not allowed to define variables only in a child context
	if secret.OriginContext.Name != config_const.DefaultContextName {

		// get all the default secrets
		defaultSecrets := c.GetSecretsByContext(config_const.DefaultContextName)
		defaultSecretFound := false

		// check if it defined
		for _, defaultSecret := range defaultSecrets {
			if defaultSecret.Name == secret.Name {
				defaultSecretFound = true
				break
			}
		}

		// return error if not defined
		if defaultSecretFound == false {
			return fmt.Errorf("secret %s defined in context %s is not defined in the default context", secret.Name, secret.OriginContext.Name)
		}

	}

	// append the secret to the repository
	c.secrets = append(c.secrets, secret)

	// sort the secrets alphabetically
	sort.SliceStable(c.secrets, func(i, j int) bool {
		return c.secrets[i].Name < c.secrets[j].Name
	})

	return nil
}

// GetSecretsByContext returns all the secrets related to the current context
func (c *Repository) GetSecretsByContext(contextName string) (res []*Secret) {
	for _, secret := range c.secrets {
		if secret.OriginContext.Name == contextName {
			res = append(res, secret)
		}
	}
	return res
}

// GetCurrentSecrets merges the default secrets with the context secrets
// the default secrets are overwritten by the context secrets
func (c *Repository) GetCurrentSecrets() (res []*Secret) {

	// get all default secrets
	defaultSecrets := c.GetSecretsByContext(config_const.DefaultContextName)

	// if not default, merge the secrets with the default secrets
	if !c.IsDefault() {

		// result is context secrets
		contextSecrets := c.GetSecretsByContext(c.context.Name)

		// add the default secret if it is missing in the context secrets
		for _, defaultSecret := range defaultSecrets {

			found := false
			for _, contextSecret := range contextSecrets {
				if defaultSecret.Name == contextSecret.Name {
					found = true
					break
				}
			}

			if found == false {
				contextSecrets = append(contextSecrets, defaultSecret)
			}

		}

		res = contextSecrets

	} else {
		// result is just the default secrets
		res = defaultSecrets
	}

	return res

}

// GetCurrentSecret takes the merged secrets from GetCurrentSecrets and returns the needed one
func (c *Repository) GetCurrentSecret(secretName string) *Secret {
	for _, secret := range c.GetCurrentSecrets() {
		if secret.Name == secretName {
			return secret
		}
	}
	return nil
}

func (s *Secret) Decode() (string, error) {
	return s.OriginContext.DecodeValue(s.EncodedValue)
}
