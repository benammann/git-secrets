package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "A brief description of your command",
	Args:  cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			value, _ := cmd.Flags().GetString("value")
			if value == "" {
				cobra.CheckErr(fmt.Errorf("you must specify a value if you dont specify a secretName"))
				cmd.Help()
			}
			decoded, errDecode := selectedContext.DecodeValue(value)
			cobra.CheckErr(errDecode)
			fmt.Println(decoded)
		} else {
			secretName := args[0]
			selectedSecret := projectCfg.GetCurrentSecret(secretName)
			if selectedSecret == nil {
				cobra.CheckErr(fmt.Errorf("the secret %s does not exist", secretName))
				return
			}

			decodedValue, errDecode := selectedSecret.Decode()
			cobra.CheckErr(errDecode)
			fmt.Println(decodedValue)
		}
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)

	decodeCmd.Flags().StringP("value", "v", "", "Decode this value instead of args[0] (key)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
