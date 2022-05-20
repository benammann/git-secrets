package cmd

import (
	cli_config "github.com/benammann/git-secrets/pkg/config/cli"
	config_const "github.com/benammann/git-secrets/pkg/config/const"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	config_parser "github.com/benammann/git-secrets/pkg/config/parser"
	"github.com/benammann/git-secrets/pkg/render"
	"github.com/spf13/cobra"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var globalCfgFile string
var projectCfgFile string
var projectCfg *config_generic.Repository
var projectCfgError error
var renderingEngine *render.RenderingEngine

var selectedContext *config_generic.Context
var contextName string

var overwrittenSecrets []string

const FlagValue = "value"
const FlagForce = "force"
const FlagDebug = "debug"
const FlagDryRun = "dry-run"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "git-secrets",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		if projectCfg != nil {
			infoCmd.Run(cmd, args)
		} else {
			cmd.Help()
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initGlobalConfig, initProjectConfig, resolveContext, createRenderingEngine)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&globalCfgFile, "global-config", "", "global config file (default is $HOME/.git-secrets.yaml)")
	rootCmd.PersistentFlags().StringVarP(&projectCfgFile, "project-config", "f", ".git-secrets.json", "project config file (default is .git-secrets.json)")
	rootCmd.PersistentFlags().StringVarP(&contextName, "context-name", "c", "", "context name (default is 'default')")
	rootCmd.PersistentFlags().StringArrayVar(&overwrittenSecrets, "secret", []string{}, "--secret secretA=$(SECRET_A_VALUE) --secret secretB=$(SECRET_B_VALUE): Pass 1-n secret names. Make sure to use environment variables to fill them!")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initGlobalConfig reads in config file and ENV variables if set.
func initGlobalConfig() {
	if globalCfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(globalCfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".git-secrets" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".git-secrets")
	}

	viper.AutomaticEnv() // read in environment variables that match
	cli_config.SetDefaults()
	viper.WatchConfig()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {

	}
}

func initProjectConfig() {

	overwrittenSecretsMap := make(map[string]string)
	for _, secretKeyValue := range overwrittenSecrets {
		splitSecret := strings.SplitN(secretKeyValue, "=", 2)
		if len(splitSecret) < 2 {
			cobra.CheckErr("Invalid Secret passed. Usage: --secret mySecret=mySecretValue")
		}
		secretKey, secretValues := splitSecret[0], splitSecret[1:]
		overwrittenSecretsMap[secretKey] = strings.Join(secretValues, "")
	}

	projectCfg, projectCfgError = config_parser.ParseRepository(projectCfgFile, overwrittenSecretsMap)

}

func resolveContext() {
	if projectCfgError != nil {
		return
	}
	desiredContextName := config_const.DefaultContextName
	if contextName != "" {
		desiredContextName = contextName
	}
	desiredContext, errGetContext := projectCfg.SetSelectedContext(desiredContextName)
	cobra.CheckErr(errGetContext)
	selectedContext = desiredContext
}

func createRenderingEngine() {
	if projectCfg != nil {
		renderingEngine = render.NewRenderingEngine(projectCfg)
	}
}
