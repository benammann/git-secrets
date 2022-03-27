package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a .git-secrets.yaml file",
	Example: `
init
init path/to/my/file.yaml
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if projectCfg != nil {
			cobra.CheckErr(fmt.Errorf("can not initialize while having config file %s loaded", projectCfgFile))
		}
		outputFile := ".git-secrets.yaml"
		if len(args) == 1 {
			outputFile = args[0]
		}

		if !strings.HasSuffix(outputFile, ".yaml") {
			cobra.CheckErr(fmt.Errorf("output file %s must have .yaml file ending", outputFile))
		}

		fmt.Println("Writing", outputFile, "...")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
