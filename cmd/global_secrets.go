package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	cli_config "github.com/benammann/git-secrets/pkg/config/cli"
	"github.com/benammann/git-secrets/pkg/encryption"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

// getGlobalSecretsCmd represents the globalSecrets command
var getGlobalSecretsCmd = &cobra.Command{
	Use:   "global-secret",
	Short: "Get or list a secret from the global configuration",
	Example: `
git-secrets get global-secrets: get all global secret keys
git-secrets get global-secret <secretKey>: prints the global secret value
`,
	Aliases: []string{"global-secrets", "gs"},
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
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

// setGlobalSecretsCmd represents the globalSecrets command
var setGlobalSecretsCmd = &cobra.Command{
	Use:   "global-secret",
	Short: "Write a secret to the global configuration",
	Example: `
git-secrets set global-secret <secretKey>: sets the global secret from terminal input
git-secrets set global-secret <secretKey> --value $MY_SECRET_VALUE_STORED_IN_ENV: sets the global secret value from --value parameter (insecure)
`,
	Aliases: []string{"global-secrets", "gs"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		isAlpha := regexp.MustCompile(`^[A-Za-z1-9]+$`).MatchString

		isForce, _ := cmd.Flags().GetBool(FlagForce)

		secretName := args[0]
		secretValue, _ := cmd.Flags().GetString(FlagValue)

		if secretValue == "" {
			errAsk := survey.AskOne(&survey.Password{
				Message: "Secret Value",
			}, &secretValue)
			cobra.CheckErr(errAsk)
		}

		if !isAlpha(secretName) {
			cobra.CheckErr("only alphanumeric letters allowed [A-Za-z1-9] allowed")
		}

		if isInvalid := encryption.ValidateAESSecret(secretValue); isInvalid != nil {
			cobra.CheckErr(fmt.Errorf("%s. hint: use 'git-secrets global-secret <key> --value $(pwgen -c 32 -n -s -y)' to generate a key", isInvalid.Error()))
		}

		finalKey := cli_config.NamedSecret(secretName)
		resolvedSecret := viper.GetString(finalKey)
		if resolvedSecret != "" && !isForce {
			cobra.CheckErr(fmt.Errorf("the secret %s already exists. Use --force to overwrite the existing secret", secretName))
		}
		viper.Set(finalKey, secretValue)
		errWrite := cli_config.WriteConfig()
		if errWrite != nil {
			cobra.CheckErr(fmt.Errorf("could not write config: %s", errWrite.Error()))
		}
		fmt.Println(finalKey, "written")
	},
}

func init() {
	getCmd.AddCommand(getGlobalSecretsCmd)
	setCmd.AddCommand(setGlobalSecretsCmd)

	setGlobalSecretsCmd.Flags().Bool(FlagForce, false, "allows to overwrite secrets")
	setGlobalSecretsCmd.Flags().String(FlagValue, "", "allows to pass the secret to write using a parameter")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getGlobalSecretsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getGlobalSecretsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
