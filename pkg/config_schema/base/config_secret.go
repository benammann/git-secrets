package base

type Secret struct {
	Name string
	OriginContext *Context
	EncodedValue string
}

func (s *Secret) Decode() (string, error) {
	return s.OriginContext.DecodeValue(s.EncodedValue)
}

func (c *Context) GetSecret(secretName string) *Secret {
	for _, secret := range c.Secrets {
		if secret.Name == secretName {
			return secret
		}
	}
	return nil
}