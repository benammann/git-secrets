package encryption

// Engine cares about encoding and decoding secrets
type Engine interface {
	EncodeValue(plainValue string) (encodedValue string, err error)
	DecodeValue(encodedValue string) (decodedValue string, err error)
}
