package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"path/filepath"
)

// setGcpCredentialsCmd represents the globalSecrets command
var setGcpCredentialsCmd = &cobra.Command{
	Use:   "gcp-credentials",
	Short: "Write a secret to the global configuration",
	Example: `
git secrets set global-secret <secretKey>: sets the global secret from terminal input
git secrets set global-secret <secretKey> --value $MY_SECRET_VALUE_STORED_IN_ENV: sets the global secret value from --value parameter (insecure)
`,
	Aliases: []string{"global-secrets", "gs"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		isForce, _ := cmd.Flags().GetBool(FlagForce)
		secretKey := args[0]
		pathToFile := args[1]

		absFilePath, errAbs := filepath.Abs(pathToFile)
		cobra.CheckErr(errAbs)

		errWrite := globalCfg.SetGcpCredentials(secretKey, absFilePath, isForce)
		if errWrite != nil {
			cobra.CheckErr(fmt.Errorf("could not write config: %s", errWrite.Error()))
		}

		fmt.Printf("%s written.\n", secretKey)

	},
}

func init() {
	setCmd.AddCommand(setGcpCredentialsCmd)

	setGcpCredentialsCmd.Flags().Bool(FlagForce, false, "Force overwrite existing secret: You may loose your master password!")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getGlobalSecretsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getGlobalSecretsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

