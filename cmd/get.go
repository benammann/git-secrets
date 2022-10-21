package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources like config, secret or global-secret",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// getConfigCmd represents the getConfig command
var getConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Get a config entry from the config file",
	Example: `
git secrets get config <configKey>
git secrets get config <configKey> -c prod
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configKey := args[0]
		configEntry := projectCfg.GetCurrentConfig(configKey)
		if configEntry == nil {
			cobra.CheckErr(fmt.Errorf("the config entry %s does not exist on context %s", configKey, selectedContext.Name))
		}
		fmt.Println(configEntry.Value)
	},
}

// getSecretCmd represents the getSecret command
var getSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Get and decode a secret entry from the config file",
	Example: `
git secrets get secret <secretName>
git secrets get secret <secretName> -c prod
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretKey := args[0]
		secretEntry := projectCfg.GetCurrentSecret(secretKey)
		if secretEntry == nil {
			cobra.CheckErr(fmt.Errorf("the secret %s does not exist on context %s", secretKey, selectedContext.Name))
		}
		decodedValue, errDecode := secretEntry.GetPlainValue(cmd.Context())
		if errDecode != nil {
			cobra.CheckErr(fmt.Errorf("could not decode secret %s: %s", secretKey, errDecode.Error()))
		}
		fmt.Println(decodedValue)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getConfigCmd)
	getCmd.AddCommand(getSecretCmd)
}
