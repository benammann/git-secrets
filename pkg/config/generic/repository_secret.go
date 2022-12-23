package config_generic

import (
	"context"
	"fmt"
	config_const "github.com/benammann/git-secrets/pkg/config/const"
	"sort"
)

type Secret interface {
	GetType() string
	GetName() string
	GetOriginContext() *Context
	GetPlainValue(ctx context.Context) (string, error)
}

// AddSecret adds a secret to the repository
// also does some validations
func (c *Repository) AddSecret(secret Secret) error {

	secretName, secretContext := secret.GetName(), secret.GetOriginContext()

	// if not default secret we need to check if the given secret is also configured in the default context
	// because we are not allowed to define variables only in a child context
	if secretContext.Name != config_const.DefaultContextName {

		// get all the default secrets
		defaultSecrets := c.GetSecretsByContext(config_const.DefaultContextName)
		defaultSecretFound := false

		// check if it defined
		for _, defaultSecret := range defaultSecrets {
			if defaultSecret.GetName() == secretName {
				defaultSecretFound = true
				break
			}
		}

		// return error if not defined
		if defaultSecretFound == false {
			return fmt.Errorf("secret %s defined in context %s is not defined in the default context", secretName, secretContext.Name)
		}

	}

	// append the secret to the repository
	c.secrets = append(c.secrets, secret)

	// sort the secrets alphabetically
	sort.SliceStable(c.secrets, func(i, j int) bool {
		return c.secrets[i].GetName() < c.secrets[j].GetName()
	})

	return nil
}

// GetSecretsByContext returns all the secrets related to the current context
func (c *Repository) GetSecretsByContext(contextName string) (res []Secret) {
	for _, secret := range c.secrets {
		if secret.GetOriginContext().Name == contextName {
			res = append(res, secret)
		}
	}
	return res
}

// GetCurrentSecrets merges the default secrets with the context secrets
// the default secrets are overwritten by the context secrets
func (c *Repository) GetCurrentSecrets() (res []Secret) {

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
				if defaultSecret.GetName() == contextSecret.GetName() {
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
func (c *Repository) GetCurrentSecret(secretName string) Secret {
	for _, secret := range c.GetCurrentSecrets() {
		if secret.GetName() == secretName {
			return secret
		}
	}
	return nil
}

// GetSecretsMapDecoded decodes the secrets of the current context and puts them into a map[string]string
func (c *Repository) GetSecretsMapDecoded(ctx context.Context) (SecretsMap, error) {

	// create the secrets map
	secretsMap := make(SecretsMap)

	// decode each secret
	for _, secret := range c.GetCurrentSecrets() {
		decodedSecret, errDecode := secret.GetPlainValue(ctx)
		if errDecode != nil {
			return nil, fmt.Errorf("could not decode secret %s: %s", secret.GetName(), errDecode.Error())
		}
		secretsMap[secret.GetName()] = decodedSecret
	}

	return secretsMap, nil

}
