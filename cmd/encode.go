package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"

	"github.com/spf13/cobra"
)

const FlagWrite = "write"
const FlagValue = "value"

// encodeCmd represents the encode command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Args:  cobra.ExactArgs(0),
	Short: "Encodes a value using the global secret configured in the current context",
	Long:  "Encodes a value using the global secret configured in the current context. Use --write <secretKey> to write the encoded value directly to the .git-secrets.json file. Use --value <value> instead of hidden input.",
	Example: `
git-secrets encode: Encodes the secret using interactive ui
git-secrets encode --write testKey: Encodes the secret and writes it to the current .git-secrets.json file
git-secrets encode --write testKey --context prod: Writes the secret to the prod context instead
git-secrets encode --write testKey --value "My Secret Value": INSECURE: Uses the value directly from the --value parameter
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Run: func(cmd *cobra.Command, args []string) {

		writeTo, _ := cmd.Flags().GetString(FlagWrite)
		value, _ := cmd.Flags().GetString(FlagValue)

		if value == "" {
			errAsk := survey.AskOne(&survey.Password{
				Message: "Value to encode",
			}, &value)
			cobra.CheckErr(errAsk)
		}

		encodedValue, errEncode := selectedContext.EncodeValue(value)
		cobra.CheckErr(errEncode)

		if writeTo != "" {
			writer := projectCfg.GetConfigWriter()
			errWrite := writer.AddSecret(projectCfg.GetCurrent().Name, writeTo, encodedValue)
			cobra.CheckErr(errWrite)
			fmt.Println("Secret", writeTo, "written to .git-secrets.json")
			fmt.Println("Get the decoded value: git-secrets decode", writeTo)

		} else {
			fmt.Println(encodedValue)
		}

	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	encodeCmd.Flags().String(FlagWrite, "", "--write <my-secret-name>: writes the secret to the current .git-secrets.json")
	encodeCmd.Flags().String(FlagValue, "", "--value <my-secret-value>: INSECURE! Uses this value instead of interactive input")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
