package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// getGlobalSecretsCmd represents the globalSecrets command
var getGlobalSecretsCmd = &cobra.Command{
	Use:   "global-secret",
	Short: "Get or list a secret from the global configuration",
	Example: `
git secrets get global-secrets: get all global secret keys
git secrets get global-secret <secretKey>: prints the global secret value
`,
	Aliases: []string{"global-secrets", "gs"},
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			secretName := args[0]
			resolvedSecret := globalCfg.GetSecret(secretName)
			if resolvedSecret != "" {
				fmt.Println(resolvedSecret)
			} else {
				cobra.CheckErr(fmt.Errorf("the secret %s does not exist", secretName))
			}
		} else {
			secretKeys := globalCfg.GetSecretKeys()
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
git secrets set global-secret <secretKey>: sets the global secret from terminal input
git secrets set global-secret <secretKey> --value $MY_SECRET_VALUE_STORED_IN_ENV: sets the global secret value from --value parameter (insecure)
`,
	Aliases: []string{"global-secrets", "gs"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		isForce, _ := cmd.Flags().GetBool(FlagForce)
		secretKey := args[0]
		secretValue, _ := cmd.Flags().GetString(FlagValue)

		if secretValue == "" {
			errAsk := survey.AskOne(&survey.Password{
				Message: "Secret Value",
			}, &secretValue)
			cobra.CheckErr(errAsk)
		}

		errWrite := globalCfg.SetSecret(secretKey, secretValue, isForce)
		if errWrite != nil {
			cobra.CheckErr(fmt.Errorf("could not write config: %s", errWrite.Error()))
		}

		fmt.Printf("%s written. Use git secrets get global-secret %s to get it's value\n", secretKey, secretKey)

	},
}

func init() {
	getCmd.AddCommand(getGlobalSecretsCmd)
	setCmd.AddCommand(setGlobalSecretsCmd)

	setGlobalSecretsCmd.Flags().Bool(FlagForce, false, "Force overwrite existing secret: You may loose your master password!")
	setGlobalSecretsCmd.Flags().String(FlagValue, "", "Pass the secret's value as parameter instead of password input")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getGlobalSecretsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getGlobalSecretsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
