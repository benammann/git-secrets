package cmd

import (
	"fmt"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	"github.com/spf13/cobra"
	"html/template"
	"os"
)

const FlagDryRun = "dry-run"

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

		for _, fileToRender := range selectedContext.FilesToRender {

			renderContext := &RenderFileData{
				UsedConfig:  projectCfgFile,
				UsedFile:    fileToRender,
				UsedContext: selectedContext.Name,
				Secrets:     decodedSecrets,
			}

			if isDryRun {
				fmt.Println("Would render file", fileToRender.FileIn, "to", fileToRender.FileOut)
			}

			tpl := template.New(fileToRender.FileIn)
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
				fmt.Println()
			} else {
				fmt.Println("File Written:", fileToRender.FileOut)
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().Bool(FlagDryRun, false, "Render files to os.stdout: --dry-run instead of writing")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
