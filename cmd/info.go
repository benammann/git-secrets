package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const InfoCmdFlagDecode = "decode"

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get detailed information about the current configuration",
	Example: `
git secrets info
git secrets info -d: Also decodes all secrets
git secrets info -d -c prod: Decodes all secrets from the prod context
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Run: func(cmd *cobra.Command, args []string) {

		allContexts := projectCfg.GetContexts()
		var allContextNames []string
		for _, context := range allContexts {
			allContextNames = append(allContextNames, context.Name)
		}

		fmt.Printf("Config File: %s (Version: %d)\n", projectCfgFile, projectCfg.GetConfigVersion())
		fmt.Printf("Available Contexts: %s\n", strings.Join(allContextNames, ", "))
		fmt.Printf("Available Render Targets: %s\n", strings.Join(projectCfg.RenderTargetNames(), ", "))
		fmt.Printf("\n")

		configHeader := []string{"Config Key", "Config Value", "Origin Context"}

		var configData [][]string

		for _, config := range projectCfg.GetCurrentConfigs() {

			tableRow := []string{config.Name, config.Value, config.OriginContext.Name}

			configData = append(configData, tableRow)

		}

		if len(configData) > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(configHeader)
			table.SetBorder(false)
			table.AppendBulk(configData)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.Render()
			fmt.Println()

		}

		shouldDecode, _ := cmd.Flags().GetBool(InfoCmdFlagDecode)

		tableHeader := []string{"Type", "Secret Name", "Origin Context"}
		if shouldDecode {
			tableHeader = append(tableHeader, "Decoded Value")
		}

		var tableData [][]string

		for _, secret := range projectCfg.GetCurrentSecrets() {

			tableRow := []string{secret.GetType(), secret.GetName(), secret.GetOriginContext().Name}
			if shouldDecode {
				decodedValue, errDecode := secret.GetPlainValue(cmd.Context())
				if errDecode != nil {
					fmt.Printf("Could not decode %s: %s\n", secret.GetName(), errDecode.Error())
					continue
				}
				tableRow = append(tableRow, decodedValue)
			}

			tableData = append(tableData, tableRow)

		}

		if len(tableData) > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(tableHeader)
			table.SetBorder(false)
			table.AppendBulk(tableData)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.Render()
			fmt.Println()

			if shouldDecode == false {
				fmt.Println("Use --decode or -d to show the decoded secrets")
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().BoolP(InfoCmdFlagDecode, "d", false, "Adds the decoded secrets to the info table")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
