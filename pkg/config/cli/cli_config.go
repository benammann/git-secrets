package cli_config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	Secrets       = "secrets"
	DaemonWatches = "daemon.watches"
)

func SetDefaults() {
	viper.SetDefault(Secrets, make(map[string]string))
	viper.SetDefault(DaemonWatches, make(map[string]string))
}

func NamedSecret(secretName string) string {
	return fmt.Sprintf("%s.%s", Secrets, secretName)
}

func WriteConfig() error {
	if _, err := os.Stat(viper.ConfigFileUsed()); os.IsNotExist(err) {
		errWrite := viper.SafeWriteConfig()
		if errWrite != nil {
			return fmt.Errorf("could not write config: %s", errWrite.Error())
		}
	} else {
		errWrite := viper.WriteConfig()
		if errWrite != nil {
			return fmt.Errorf("could not write config: %s", errWrite.Error())
		}
	}
	return nil
}
