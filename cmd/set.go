package cmd

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/benammann/git-secrets/pkg/gcp"
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

var setSecretCmd = &cobra.Command{
	Use: "secret",
	Short: "Commands to write secrets to the config",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// setEncryptedSecretCmd represents the setSecret command
var setEncryptedSecretCmd = &cobra.Command{
	Use:   "encrypted",
	Short: "Encode and write a secret to the config file",
	Example: `
git secrets set secret encrypted <secretKey>: Encodes the secret using interactive ui and adds it to the git secrets file
git secrets set secret encrypted <secretKey> --value <plainValue>: INSECURE: Uses the value directly from the --value parameter
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

// setEncryptedSecretCmd represents the setSecret command
var setGcpSecretCommand = &cobra.Command{
	Use:   "gcp",
	Short: "Encode and write a secret to the config file",
	Example: `
git secrets set secret gcp <secretKey>: Resolves the available secrets from secret manager
git secrets set secret gcp <secretKey> --resourceId <resourceId>: Uses the resource url directly from the --resourceId parameter.
`,
	Args: cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		secretKey := args[0]
		resourceId, _ := cmd.Flags().GetString(FlagResourceId)
		force, _ := cmd.Flags().GetBool(FlagForce)

		if projectCfg.GetDefault().GcpCredentials == "" {
			return fmt.Errorf("you need to configure gcp credentials first")
		}
		credentialsName, errCredentials := projectCfg.GetCurrentGCPCredentialsName()
		cobra.CheckErr(errCredentials)
		credentialsFileName := globalCfg.GetGcpCredentialsFile(credentialsName)

		if resourceId == "" {
			gcpSecrets, errSecrets := gcp.ListSecrets(cmd.Context(), credentialsFileName)
			cobra.CheckErr(errSecrets)
			for _, secret := range gcpSecrets {
				fmt.Println(secret.GetEtag())
			}
			cobra.CheckErr(errors.New("selection not supported yet. use --resourceId <resourceId>"))
		}

		writer := projectCfg.GetConfigWriter()

		errWrite := writer.SetGcpSecret(projectCfg.GetCurrent().Name, secretKey, resourceId, force)
		cobra.CheckErr(errWrite)

		fmt.Printf("The secret %s has been written\n", secretKey)
		fmt.Printf("Resolve the decoded value: git secrets get secret %s\n", secretKey)
		fmt.Printf("Use it in a template: MY_CONFIG_KEY={{.Secrets.%s}}\n", secretKey)

		return nil

	},
}

func init() {
	for _, cmd := range []*cobra.Command{setConfigCmd, setEncryptedSecretCmd, setGcpSecretCommand} {
		cmd.Flags().Bool(FlagForce, false, "use --force to overwrite an existing value")
	}
	setEncryptedSecretCmd.Flags().String(FlagValue, "", "--value <secretValue>: This is insecure, use --value $ENV_VALUE to not write the secret value to the history file.")
	setGcpSecretCommand.Flags().String(FlagResourceId, "", "--resourceId <resourceId>: Uses the resource url directly from the --resourceId parameter.")
	rootCmd.AddCommand(setCmd)

	setSecretCmd.AddCommand(setEncryptedSecretCmd)
	setSecretCmd.AddCommand(setGcpSecretCommand)

	setCmd.AddCommand(setConfigCmd)
	setCmd.AddCommand(setSecretCmd)
}
