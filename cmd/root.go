/*
Copyright © 2022 Ben Ammann

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/benammann/git-secrets/pkg/config"
	"github.com/benammann/git-secrets/pkg/encryption"
	"github.com/spf13/cobra"
	"log"
	"os"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-secrets",
	Args: cobra.ExactArgs(1),
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]
		configOut, errConfig := config.ParseConfig(fileName)
		cobra.CheckErr(errConfig)
		defaultContext := configOut.GetDefaultContext()
		encryptionEngine := encryption.NewAesEngine(defaultContext.SecretResolver)

		encrypted, errEncrypt := encryptionEngine.EncodeValue("Hallo Welt")
		if errEncrypt != nil {
			log.Fatalf("could not encrypt: %s", errEncrypt.Error())
		}

		decrypted, errDecrypt := encryptionEngine.DecodeValue(encrypted)
		if errDecrypt != nil {
			log.Fatalf("could not decrypt: %s", errDecrypt.Error())
		}

		encryptedBase64 := base64.StdEncoding.EncodeToString([]byte(encrypted))

		fmt.Println("Encrypted:", encryptedBase64)
		fmt.Println("Decrypted:", decrypted)

	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.git-secrets.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".git-secrets" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".git-secrets")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
