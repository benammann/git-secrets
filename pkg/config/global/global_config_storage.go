package global_config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type StorageProvider interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	GetString(key string) string
	AllKeys() []string
	WriteConfig() error
}

type ViperStorageProvider struct {
	viperInstance *viper.Viper
}

type MemoryStorageProvider struct {
	storage map[string]interface{}
}

func NewViperConfigStorage(viperInstance *viper.Viper) *ViperStorageProvider {
	return &ViperStorageProvider{
		viperInstance: viperInstance,
	}
}

func NewMemoryStorageProvider() *MemoryStorageProvider {
	return &MemoryStorageProvider{
		storage: make(map[string]interface{}),
	}
}

func (v *ViperStorageProvider) Set(key string, value interface{}) {
	v.viperInstance.Set(key, value)
}

func (v *ViperStorageProvider) Get(key string) interface{} {
	return v.viperInstance.Get(key)
}

func (v *ViperStorageProvider) GetString(key string) string {
	return v.viperInstance.GetString(key)
}

func (v *ViperStorageProvider) AllKeys() []string {
	return v.viperInstance.AllKeys()
}

func (v *ViperStorageProvider) WriteConfig() error {
	if _, err := os.Stat(v.viperInstance.ConfigFileUsed()); os.IsNotExist(err) {
		errWrite := v.viperInstance.SafeWriteConfig()
		if errWrite != nil {
			return fmt.Errorf("could not write config: %s", errWrite.Error())
		}
	} else {
		errWrite := v.viperInstance.WriteConfig()
		if errWrite != nil {
			return fmt.Errorf("could not write config: %s", errWrite.Error())
		}
	}
	return nil
}

func (m *MemoryStorageProvider) Set(key string, value interface{}) {
	m.storage[key] = value
}

func (m *MemoryStorageProvider) Get(key string) interface{} {
	return m.storage[key]
}

func (m *MemoryStorageProvider) GetString(key string) string {
	value := m.Get(key)
	switch s := value.(type) {
	case string:
		return s
	default:
		return ""
	}
}

func (m *MemoryStorageProvider) AllKeys() []string {
	var keys []string
	for key, _ := range m.storage {
		keys = append(keys, key)
	}
	return keys
}

func (m *MemoryStorageProvider) WriteConfig() error {
	return nil
}
