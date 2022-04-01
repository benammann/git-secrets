package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	"github.com/spf13/cobra"
	"html/template"
)

const FlagDebug = "debug"
const FlagDryRun = "dry-run"
const FlagFileIn = "file-in"
const FlagFileOut = "file-out"

type RenderFileData struct {
	UsedConfig  string
	UsedContext string
	UsedFile    *config_generic.FileToRender
	Secrets     map[string]string
}

func Base64Encode(args ...interface{}) string {
	return base64.StdEncoding.EncodeToString([]byte(args[0].(string)))
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"Base64Encode": Base64Encode,
	}
}

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "render files feature",
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
	renderCmd.Flags().StringP(FlagFileIn, "i", "", "Input file to render (requires also --file-out or -o flag)")
	renderCmd.Flags().StringP(FlagFileOut, "o", "", "Output file to render (requires also --file-in or -i flag)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
