package cmd

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	config_init "github.com/benammann/git-secrets/pkg/config/init"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new .git-secrets.json project",
	Example: `
git secrets init
git secrets init -f custom-file-name.json
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if projectCfg != nil {
			cobra.CheckErr(fmt.Errorf("can not initialize while having config file %s loaded. Please switch directories", projectCfgFile))
		}

		secretKeys := globalCfg.GetSecretKeys()

		if len(secretKeys) < 0 {
			cobra.CheckErr(fmt.Errorf("please create a global secret before: git secrets set global-secret <secret-name>"))
		}

		var outputFileQuestions = []*survey.Question{
			{
				Name: "outputFile",
				Prompt: &survey.Input{
					Message: "Output file",
					Default: ".git-secrets.json",
				},
				Validate: func(ans interface{}) error {

					outputFile := ans.(string)
					if outputFile == "" {
						return fmt.Errorf("the output file cannot be empty")
					}

					if _, err := os.Stat(outputFile); errors.Is(err, os.ErrExist) {
						return fmt.Errorf("%s already exists", outputFile)
					}

					return nil

				},
			},
			{
				Name: "secretName",
				Prompt: &survey.Select{
					Message: "Which global secret should be used to encode / decode secrets:",
					Options: secretKeys,
				},
				Validate: survey.Required,
			},
		}

		questionResponse := struct {
			OutputFile string `survey:"outputFile"`
			SecretName string `survey:"secretName"`
		}{}

		if errAsk := survey.Ask(outputFileQuestions, &questionResponse); errAsk != nil {
			cobra.CheckErr(fmt.Errorf("could not ask survey: %s", errAsk.Error()))
		}

		if !strings.HasSuffix(questionResponse.OutputFile, ".json") {
			cobra.CheckErr(fmt.Errorf("output file %s must have .json file ending", questionResponse.OutputFile))
		}

		errWrite := config_init.WriteInitialConfig(fs, questionResponse.OutputFile, questionResponse.SecretName)
		if errWrite != nil {
			cobra.CheckErr(errWrite)
		}

		fmt.Println(questionResponse.OutputFile, "written")
		fmt.Println("Info: git secrets info -d")
		fmt.Println("Add Context: git secrets add context <contextName>")
		fmt.Println("Set Config: git secrets set config <configKey> <configValue>")
		fmt.Println("Set Secret: git secrets set secret <secretKey>")

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
