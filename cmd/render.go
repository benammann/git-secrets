package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	"github.com/spf13/cobra"
	"html/template"
	"os"
	"path"
)

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

		decodedSecrets := make(map[string]string)
		for _, secret := range projectCfg.GetCurrentSecrets() {
			decodedSecret, errDecode := secret.Decode()
			if errDecode != nil {
				cobra.CheckErr(fmt.Errorf("could not decode secret %s: %s", secret.Name, errDecode.Error()))
			}
			decodedSecrets[secret.Name] = decodedSecret
		}

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

		fmt.Printf("Rendering as context %s ...\n\n", selectedContext.Name)

		for _, fileToRender := range filesToRender {

			renderContext := &RenderFileData{
				UsedConfig:  projectCfgFile,
				UsedFile:    fileToRender,
				UsedContext: selectedContext.Name,
				Secrets:     decodedSecrets,
			}

			if isDryRun {
				fmt.Println("")
				fmt.Println("Would render file", fileToRender.FileIn, "to", fileToRender.FileOut)
				renderContextJson, _ := json.MarshalIndent(renderContext, "", "  ")
				fmt.Println("Available Variables:")
				fmt.Println(string(renderContextJson))
				fmt.Println("")
			}

			tpl := template.New(path.Base(fileToRender.FileIn))
			tpl.Funcs(getFuncMap())
			tpl, errTpl := tpl.ParseFiles(fileToRender.FileIn)
			if errTpl != nil {
				cobra.CheckErr(fmt.Errorf("could not read file contents of %s: %s", fileToRender.FileIn, errTpl.Error()))
			}

			fileOut := os.Stdout
			if isDryRun == false {
				fsFileOut, errFsFile := os.OpenFile(fileToRender.FileOut, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
				if fsFileOut != nil {
					defer fsFileOut.Close()
				}
				if errFsFile != nil {
					cobra.CheckErr(fmt.Errorf("could not open output file %s:%s", fileToRender.FileOut, errFsFile.Error()))
				}
				fileOut = fsFileOut
			}

			errExecute := tpl.Execute(fileOut, renderContext)
			if errExecute != nil {
				cobra.CheckErr(fmt.Errorf("could not render file %s: %s", fileToRender.FileOut, errExecute.Error()))
			}

			if isDryRun {
				fmt.Println()
			} else {
				fmt.Println(fileToRender.FileIn, "->", fileToRender.FileOut)
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().Bool(FlagDryRun, false, "Render files to os.stdout: --dry-run instead of writing")
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
