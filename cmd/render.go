package cmd

import (
	"encoding/json"
	"fmt"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	"github.com/spf13/cobra"
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
	Short: "render files feature",
	Example: `
git-secrets render: Render from configuration
git-secrets render <fileIn> <fileOut> --debug: Render a specific file instead of the configured ones
git-secrets render -c prod: Render files for the prod context
git-secrets render --dry-run: Render files and print them to the console
git-secrets render --dry-run --debug: Dry run render and shows the rendering context
git-secrets render --debug: Render and write the rendering context
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if !(len(args) == 0 || len(args) == 2) {
			return fmt.Errorf("usage: git-secrets render or git-secrets render <file-in> <file-out>")
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
			if len(selectedContext.FilesToRender) == 0 {
				cobra.CheckErr(fmt.Errorf("the context %s has no files to render. Use --file-in to render a custom file using this context", selectedContext.Name))
			}
			filesToRender = selectedContext.FilesToRender
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
