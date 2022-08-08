package cmd

import (
	"encoding/json"
	"fmt"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	"github.com/spf13/cobra"
	"strings"
)

type RenderFileData struct {
	UsedConfig  string
	UsedContext string
	UsedFile    *config_generic.FileToRender
	Secrets     map[string]string
}

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render files using the go templating engine",
	Example: `
git-secrets render <targetName>: Render from configuration
git-secrets render <targetName1>,<targetName2>,...: Renders multiple targets at once
git-secrets render <fileIn> <fileOut> --debug: Render a specific file instead of the configured ones
git-secrets render <targetName> -c prod: Render files for the prod context
git-secrets render <targetName> --dry-run: Render files and print them to the console
git-secrets render <targetName> --dry-run --debug: Dry run render and shows the rendering context
git-secrets render <targetName> --debug: Render and write the rendering target
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if !(len(args) == 0 || len(args) == 1 || len(args) == 2) {
			return fmt.Errorf("usage: git-secrets render <target> or git-secrets render <file-in> <file-out>")
		}
		return nil
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Run: func(cmd *cobra.Command, args []string) {

		isDryRun, _ := cmd.Flags().GetBool(FlagDryRun)
		isDebug, _ := cmd.Flags().GetBool(FlagDebug)

		var filesToRender []*config_generic.FileToRender
		if len(args) == 0 {
			fmt.Println("Usage: git-secrets render <targetName1>,<targetName2>,...")
			fmt.Println("Render using another context: git-secrets render <targetName> -c <contextName>")
			cobra.CheckErr(fmt.Errorf("You must specify a rendering context. Available targets: %s", strings.Join(projectCfg.RenderTargetNames(), ", ")))
		} else if len(args) == 1 {

			nonUniqueTargets := strings.Split(args[0], ",")
			uniqueTargets := make(map[string]bool)

			for _, nonUniqueTarget := range nonUniqueTargets {
				if !uniqueTargets[nonUniqueTarget] {
					uniqueTargets[nonUniqueTarget] = true
				}
			}

			for requestedTargetName := range uniqueTargets {
				requestedTarget := projectCfg.GetRenderTarget(requestedTargetName)
				if requestedTarget == nil {
					cobra.CheckErr(fmt.Errorf("the render target %s does not exist. Available targets: %s", args[0], strings.Join(projectCfg.RenderTargetNames(), ", ")))
				}
				filesToRender = append(filesToRender, requestedTarget.FilesToRender...)
			}

			if len(filesToRender) == 0 {
				fmt.Println("could not resolve any files to render. Use git secrets render <fileIn> <fileOut> -c <contextName> to render a file manually")
				cobra.CheckErr(fmt.Errorf("you can also add a file using git secrets add file <fileIn> <fileOut> -t <renderTarget>"))
			}

		} else {
			filesToRender = append(filesToRender, &config_generic.FileToRender{
				FileIn:  args[0],
				FileOut: args[1],
			})
		}

		for _, fileToRender := range filesToRender {

			if isDryRun {
				usedContext, fileContents, errRender := renderingEngine.RenderFile(fileToRender)
				if isDebug {
					fmt.Println(fileToRender.FileIn)
					if usedContext != nil {
						renderContextJson, _ := json.MarshalIndent(usedContext, "", "  ")
						fmt.Println(string(renderContextJson))
					}
				}
				if errRender != nil {
					cobra.CheckErr(fmt.Errorf("could not render file %s: %s", fileToRender.FileIn, errRender.Error()))
					continue
				}
				fmt.Println(fileContents)
			} else {
				usedContext, errWrite := renderingEngine.WriteFile(fileToRender)
				if isDebug && usedContext != nil {
					fmt.Println(fileToRender.FileIn)
					renderContextJson, _ := json.MarshalIndent(usedContext, "", "  ")
					fmt.Println(string(renderContextJson))
				}
				if errWrite != nil {
					cobra.CheckErr(fmt.Errorf("could not write file %s: %s", fileToRender.FileIn, errWrite.Error()))
					continue
				}
				fmt.Println(fileToRender.FileOut, "written")
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().Bool(FlagDryRun, false, "Render files to os.stdout: --dry-run instead of writing")
	renderCmd.Flags().Bool(FlagDebug, false, "Also prints the rendering context to the console")

}
