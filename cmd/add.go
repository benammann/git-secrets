package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "allows you to set resources in your projects or global config file",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// addContextCmd represents the addContext command
var addContextCmd = &cobra.Command{
	Use:     "context",
	Short:   "adds a context to the existing config file",
	Example: "git-secrets add context <contextName>",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		writer := projectCfg.GetConfigWriter()
		contextToAdd := args[0]
		errAdd := writer.AddContext(contextToAdd)
		cobra.CheckErr(errAdd)
		fmt.Printf("The context %s has been added to your config file\n", contextToAdd)
		fmt.Printf("Now use it using the --context %s or -c %s flag\n", contextToAdd, contextToAdd)
		fmt.Printf("Add a config to this context: git-secrets set config <configKey> <configValue> -c %s\n", contextToAdd)
	},
}

// addFileCmd represents the addFile command
var addFileCmd = &cobra.Command{
	Use:   "file",
	Short: "adds a file to render to the git-secrets file",
	Example: `
git-secrets add file <fileIn> <fileOut>
git-secrets add file <fileIn> <fileOut> -c prod
`,
	Args: cobra.ExactArgs(2),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fileIn, fileOut := args[0], args[1]
		configWrite := projectCfg.GetConfigWriter()
		cobra.CheckErr(configWrite.AddFileToRender(projectCfg.GetCurrent().Name, fileIn, fileOut))
		fmt.Printf("Render File %s/%s has been added to your config file.\n", fileIn, fileOut)
		fmt.Println("To render the file use: git-secrets render")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addContextCmd)
	addCmd.AddCommand(addFileCmd)
}
