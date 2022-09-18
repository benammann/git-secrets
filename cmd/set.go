package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set resources like config, secret or global-secret",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// setConfigCmd represents the setConfig command
var setConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Set a config entry",
	Example: `
git secrets set config <configKey> <configValue>
git secrets set config <configKey> <configValue> -c prod
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool(FlagForce)
		configKey, configValue := args[0], args[1]
		configWrite := projectCfg.GetConfigWriter()
		cobra.CheckErr(configWrite.SetConfig(projectCfg.GetCurrent().Name, configKey, configValue, force))
		fmt.Printf("The config entry %s has been written\n", configKey)
		fmt.Printf("Resolve the value: git secrets get config %s\n", configKey)
		fmt.Printf("Use it in a template: MY_CONFIG_KEY={{.Configs.%s}}\n", configKey)
	},
}

// setSecretCmd represents the setSecret command
var setSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Encode and write a secret to the config file",
	Example: `
git secrets set secret <secretKey>: Encodes the secret using interactive ui and adds it to the git secrets file
git secrets set secret <secretKey> --value <plainValue>: INSECURE: Uses the value directly from the --value parameter
`,
	Args: cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Run: func(cmd *cobra.Command, args []string) {

		secretKey := args[0]
		value, _ := cmd.Flags().GetString(FlagValue)
		force, _ := cmd.Flags().GetBool(FlagForce)

		if value == "" {
			errAsk := survey.AskOne(&survey.Password{
				Message: "Value to encode",
			}, &value)
			cobra.CheckErr(errAsk)
		}

		encodedValue, errEncode := selectedContext.EncodeValue(value)
		cobra.CheckErr(errEncode)

		writer := projectCfg.GetConfigWriter()

		errWrite := writer.SetEncryptedSecret(projectCfg.GetCurrent().Name, secretKey, encodedValue, force)
		cobra.CheckErr(errWrite)

		fmt.Printf("The secret %s has been written\n", secretKey)
		fmt.Printf("Resolve the decoded value: git secrets get secret %s\n", secretKey)
		fmt.Printf("Use it in a template: MY_CONFIG_KEY={{.Secrets.%s}}\n", secretKey)

	},
}

func init() {
	for _, cmd := range []*cobra.Command{setConfigCmd, setSecretCmd} {
		cmd.Flags().Bool(FlagForce, false, "use --force to overwrite an existing value")
	}
	setSecretCmd.Flags().String(FlagValue, "", "--value <secretValue>: This is insecure, use --value $ENV_VALUE to not write the secret value to the history file.")
	rootCmd.AddCommand(setCmd)
	setCmd.AddCommand(setConfigCmd)
	setCmd.AddCommand(setSecretCmd)
}
