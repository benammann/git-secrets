package config_generic


func NewEncryptedSecret(name string, encodedValue string, originContext *Context) *EncryptedSecret {
	return &EncryptedSecret{
		Name: name,
		EncodedValue: encodedValue,
		OriginContext: originContext,
	}
}

type EncryptedSecret struct {

	// Name describes the name of the GcpSecret
	Name string

	// EncodedValue hold the encodedValue in base64 of the GcpSecret
	EncodedValue string

	// OriginContext references the configured context to decode the GcpSecret
	OriginContext *Context
}

func (s *EncryptedSecret) GetName() string {
	return s.Name
}

func (s *EncryptedSecret) GetOriginContext() *Context {
	return s.OriginContext
}

func (s *EncryptedSecret) GetPlainValue() (string, error) {
	return s.OriginContext.DecodeValue(s.EncodedValue)
}

