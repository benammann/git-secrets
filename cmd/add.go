package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add resources like context or file",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// addContextCmd represents the addContext command
var addContextCmd = &cobra.Command{
	Use:     "context",
	Short:   "Add a context to the config file",
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
		fmt.Printf("Add a secret to this context: git-secrets set secret <secretKey> -c %s\n", contextToAdd)
	},
}

// addFileCmd represents the addFile command
var addFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Add a file to the rendering engine",
	Example: `
git-secrets add file <fileIn> <fileOut>
git-secrets add file <fileIn> <fileOut> -c prod
`,
	Args: cobra.ExactArgs(2),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Run: func(cmd *cobra.Command, args []string) {
		targetName, _ := cmd.Flags().GetString(FlagTarget)
		if targetName == "" {
			cobra.CheckErr(fmt.Errorf("you must specify a target name: -t or --target <targetName>"))
		}
		fileIn, fileOut := args[0], args[1]
		configWrite := projectCfg.GetConfigWriter()
		cobra.CheckErr(configWrite.AddFileToRender(targetName, fileIn, fileOut))
		fmt.Printf("Render File %s/%s has been added to your config file.\n", fileIn, fileOut)
		fmt.Printf("To render the file use: git-secrets render %s or git-secrets render %s -c <contextName>\n", targetName, targetName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addContextCmd)
	addCmd.AddCommand(addFileCmd)
	addFileCmd.Flags().StringP(FlagTarget, "t", "", "Specifies the render target name: -t <targetName>, example -t k8s")
}
