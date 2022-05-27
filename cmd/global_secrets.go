/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	cli_config "github.com/benammann/git-secrets/pkg/config/cli"
	"github.com/benammann/git-secrets/pkg/encryption"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

// globalSecretsCmd represents the globalSecrets command
var globalSecretsCmd = &cobra.Command{
	Use:   "global-secret",
	Short: "allows to manage the global secrets from ~/.git-secrets.yaml using the cli",
	Example: `
git-secrets global-secrets: get all global secret keys
git-secrets global-secret <secretKey>: prints the global secret value
git-secrets global-secret <secretKey> <secretValue>: sets the global secret value
`,
	Aliases: []string{"global-secrets"},
	Args:    cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(viper.AllKeys())

		isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString

		isSet := len(args) == 2
		isGet := len(args) == 1
		isForce, _ := cmd.Flags().GetBool(FlagForce)

		if isSet {
			secretName := args[0]
			secretValue := args[1]

			if !isAlpha(secretName) {
				cobra.CheckErr("only alphanumeric letters allowed [A-Za-z] allowed")
			}

			if isInvalid := encryption.ValidateAESSecret(secretValue); isInvalid != nil {
				cobra.CheckErr(fmt.Errorf("%s. hint: use 'gs global-secret <key> $(pwgen -c 32 -n -s -y)' to generate a key", isInvalid.Error()))
			}

			finalKey := cli_config.NamedSecret(secretName)
			resolvedSecret := viper.GetString(finalKey)
			if resolvedSecret != "" && !isForce {
				cobra.CheckErr(fmt.Errorf("the secret %s already exists. Use --force to overwrite the existing secret", secretName))
			}
			viper.Set(finalKey, args[1])
			errWrite := cli_config.WriteConfig()
			if errWrite != nil {
				cobra.CheckErr(fmt.Errorf("could not write config: %s", errWrite.Error()))
			}
			fmt.Println(finalKey, "written")
		} else if isGet {
			secretName := args[0]
			resolvedSecret := viper.GetString(cli_config.NamedSecret(secretName))
			if resolvedSecret != "" {
				fmt.Println(resolvedSecret)
			} else {
				cobra.CheckErr(fmt.Errorf("the secret %s does not exist", secretName))
			}
		} else {
			var secretKeys []string
			for _, key := range viper.AllKeys() {
				secretPrefix := fmt.Sprintf("%s.", cli_config.Secrets)
				if strings.HasPrefix(key, secretPrefix) {
					secretKeys = append(secretKeys, strings.Replace(key, secretPrefix, "", 1))
				}
			}
			for _, secretKey := range secretKeys {
				fmt.Println(secretKey)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(globalSecretsCmd)

	globalSecretsCmd.Flags().Bool(FlagForce, false, "allows to overwrite secrets")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// globalSecretsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// globalSecretsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
