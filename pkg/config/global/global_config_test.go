package global_config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGlobalConfigProvider_GetSecret(t *testing.T) {

	globalCfg := NewGlobalConfigProvider(NewMemoryStorageProvider())

	secretKey := "testSecret"
	secretValue := "teiw1ohhahpeib7Rae9erahcaik1Mion"

	t.Run("return empty string on non existing secret", func(t *testing.T) {
		assert.Equal(t, "", globalCfg.GetSecret(secretKey))
	})

	t.Run("return the correct secret value if set", func(t *testing.T) {
		assert.Nil(t, globalCfg.SetSecret(secretKey, secretValue, true))
		assert.Equal(t, secretValue, globalCfg.GetSecret(secretKey))
	})

}

func TestGlobalConfigProvider_GetSecretKeys(t *testing.T) {

	globalCfg := NewGlobalConfigProvider(NewMemoryStorageProvider())
	secretKeys := []string{"secretKeyA", "secretKeyB"}

	for _, secretKey := range secretKeys {
		assert.Nil(t, globalCfg.SetSecret(secretKey, "teiw1ohhahpeib7Rae9erahcaik1Mion", true))
	}

	assert.Len(t, globalCfg.GetSecretKeys(), len(secretKeys))
	for _, secretKey := range secretKeys {
		assert.Contains(t, globalCfg.GetSecretKeys(), strings.ToLower(secretKey))
	}

}

func TestGlobalConfigProvider_SetSecret(t *testing.T) {

	globalCfg := NewGlobalConfigProvider(NewMemoryStorageProvider())

	secretKey := "secretKey"
	secretValue := "Zu5Ousi7phohsheewooMeex2saegiQu5"
	newSecretValue := "Ohxaenu4ih3hoongogauveiMo4kiin1o"

	t.Run("it should set the secret value", func(t *testing.T) {
		assert.Nil(t, globalCfg.SetSecret(secretKey, secretValue, false))
		assert.Equal(t, secretValue, globalCfg.GetSecret(secretKey))
	})

	t.Run("it should throw an error when setting the key twice without forcing", func(t *testing.T) {
		assert.NotNil(t, globalCfg.SetSecret(secretKey, newSecretValue, false))
		assert.Equal(t, secretValue, globalCfg.GetSecret(secretKey))
	})

	t.Run("it should allow to redeclare a key by force overwrite", func(t *testing.T) {
		assert.Nil(t, globalCfg.SetSecret(secretKey, newSecretValue, true))
		assert.Equal(t, newSecretValue, globalCfg.GetSecret(secretKey))
	})

}

func TestGlobalConfigProvider_secretConfigKey(t *testing.T) {
	globalCfg := NewGlobalConfigProvider(NewMemoryStorageProvider())

	secretKey := "mySecretKey"
	expectedKey := fmt.Sprintf("%s.%s", SecretKeyPrefix, strings.ToLower(secretKey))

	assert.Equal(t, expectedKey, globalCfg.secretConfigKey(secretKey))

}

func TestGlobalConfigProvider_validateSecretValue(t *testing.T) {

	globalCfg := NewGlobalConfigProvider(NewMemoryStorageProvider())

	type args struct {
		secretValue string
	}
	tests := []struct {
		name        string
		secretKey   string
		secretValue string
		wantErr     bool
	}{
		{
			name:        "it should accept valid secrets",
			secretKey:   "mySecretKey123",
			secretValue: "Yi2ubah5xae2fou4Chairahkeisahsh7",
			wantErr:     false,
		},
		{
			name:        "it should reject invalid secret keys",
			secretKey:   "my secret key",
			secretValue: "Yi2ubah5xae2fou4Chairahkeisahsh7",
			wantErr:     true,
		},
		{
			name:        "it should reject invalid secret keys",
			secretKey:   "my!SecretKey",
			secretValue: "Yi2ubah5xae2fou4Chairahkeisahsh7",
			wantErr:     true,
		},
		{
			name:        "it should reject invalid secret value sizes",
			secretKey:   "mySecretKey123",
			secretValue: "abc",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := globalCfg.validateSecret(tt.secretKey, tt.secretValue); (err != nil) != tt.wantErr {
				t.Errorf("validateSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateAESSecret(t *testing.T) {
	b32 := make([]byte, 32)
	b24 := make([]byte, 24)
	b16 := make([]byte, 16)
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "accept 32 bytes",
			value:   string(b32),
			wantErr: false,
		},
		{
			name:    "accept 24 bytes",
			value:   string(b24),
			wantErr: false,
		},
		{
			name:    "accept 16 bytes",
			value:   string(b16),
			wantErr: false,
		},
		{
			name:    "reject other sizes",
			value:   string([]byte{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAESSecret(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("validateAESSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
