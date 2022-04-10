package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const FlagWrite = "write"

// encodeCmd represents the encode command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Args:  cobra.MinimumNArgs(1),
	Short: "Encodes a value",
	Long:  "Encodes a value",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Run: func(cmd *cobra.Command, args []string) {

		writeTo, _ := cmd.Flags().GetString(FlagWrite)

		encodedValue, errEncode := selectedContext.EncodeValue(args[0])
		cobra.CheckErr(errEncode)

		if writeTo != "" {
			writer := projectCfg.GetConfigWriter()
			errWrite := writer.AddSecret(projectCfg.GetCurrent().Name, writeTo, encodedValue)
			cobra.CheckErr(errWrite)
			fmt.Println("Secret", writeTo, "written to .git-secrets.json")
		} else {
			fmt.Println(encodedValue)
		}

	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	encodeCmd.Flags().StringP(FlagWrite, "w", "", "-w <my-secret-name>: writes the secret to the current .git-secrets.json")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
