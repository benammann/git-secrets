package cmd

import (
	"fmt"
	config_init "github.com/benammann/git-secrets/pkg/config/init"
	"github.com/spf13/cobra"
	"strings"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a .git-secrets.json file",
	Example: `
init
init path/to/my/file.yaml
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if projectCfg != nil {
			cobra.CheckErr(fmt.Errorf("can not initialize while having config file %s loaded. Please switch directories", projectCfgFile))
		}
		outputFile := ".git-secrets.json"
		if len(args) == 1 {
			outputFile = args[0]
		}

		if !strings.HasSuffix(outputFile, ".json") {
			cobra.CheckErr(fmt.Errorf("output file %s must have .json file ending", outputFile))
		}

		errWrite := config_init.WriteInitialConfig(outputFile)
		if errWrite != nil {
			cobra.CheckErr(errWrite)
		}

		fmt.Println(outputFile, "written")

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
