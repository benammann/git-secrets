package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// Engine cares about encoding and decoding secrets
type Engine interface {
	EncodeValue(plainValue string) (encodedValue string, err error)
	DecodeValue(encodedValue string) (decodedValue string, err error)
}

type AesEngine struct {
	secretResolver SecretResolver
}

func NewAesEngine(secretResolver SecretResolver) Engine {
	return &AesEngine{
		secretResolver: secretResolver,
	}
}

func (a *AesEngine) newGcm() ([]byte, cipher.AEAD, error) {

	// resolve the secret from the abstract secret resolver
	secret, errSecret := a.secretResolver.GetPlainSecret()
	if errSecret != nil {
		return nil, nil, fmt.Errorf("could not resolve secret: %s", errSecret.Error())
	}

	// create the new cipher
	newCipher, errCipher := aes.NewCipher(secret)
	if errCipher != nil {
		return nil, nil, fmt.Errorf("could not create cipher instance from secret: %s", errCipher.Error())
	}

	// create a gcm instance from cipher instance
	newGcm, errGcm := cipher.NewGCM(newCipher)
	if errGcm != nil {
		return nil, nil, fmt.Errorf("could not create gcm from cipher instance: %s", errGcm.Error())
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, newGcm.NonceSize())

	// populates our nonce with a cryptographically secure
	// random sequence
	if _, errCreateNonce := io.ReadFull(rand.Reader, nonce); errCreateNonce != nil {
		return nil, nil, fmt.Errorf("could not create nonce from gcm instance: %s", errCreateNonce.Error())
	}

	// return the new gcm
	return nonce, newGcm, nil

}

func (a *AesEngine) EncodeValue(plainValue string) (encodedValue string, err error) {
	nonce, gcm, errGcm := a.newGcm()
	if errGcm != nil {
		return "", errGcm
	}
	return string(gcm.Seal(nonce, nonce, []byte(plainValue), nil)), nil
}

func (a *AesEngine) DecodeValue(encodedValue string) (decodedValue string, err error) {

	_, gcm, errGcm := a.newGcm()
	if errGcm != nil {
		return "", errGcm
	}
	nonceSize := gcm.NonceSize()

	encodedValueBytes := []byte(encodedValue)
	if len(encodedValueBytes) < nonceSize {
		return "", fmt.Errorf("encoded value is smaller than nonce size")
	}

	nonce, cipherText := encodedValueBytes[:nonceSize], encodedValueBytes[nonceSize:]
	plainBytes, errOpen := gcm.Open(nil, nonce, cipherText, nil)
	if errOpen != nil {
		return "", fmt.Errorf("could not open via gcm: %s", errOpen.Error())
	}

	return string(plainBytes), nil

}
