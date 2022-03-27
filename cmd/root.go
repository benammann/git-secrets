package cmd

import (
	config_const "github.com/benammann/git-secrets/pkg/config/const"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	config_parser "github.com/benammann/git-secrets/pkg/config/parser"
	"github.com/spf13/cobra"
	"os"

	"github.com/spf13/viper"
)

var globalCfgFile string
var projectCfgFile string
var projectCfg *config_generic.Repository
var projectCfgError error

var selectedContext *config_generic.Context
var contextName string

var overwrites config_generic.ConfigCliArgs

var overwriteSecret string
var overwriteSecretName string
var overwriteSecretEnv string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-secrets",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	cobra.OnInitialize(initGlobalConfig, initProjectConfig, resolveContext)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&globalCfgFile, "global-config", "", "global config file (default is $HOME/.git-secrets.yaml)")
	rootCmd.PersistentFlags().StringVarP(&projectCfgFile, "project-config", "f", ".git-secrets.yaml", "project config file (default is .git-secrets.yaml)")
	rootCmd.PersistentFlags().StringVarP(&contextName, "context-name", "c", "", "context name (default is 'default')")
	rootCmd.PersistentFlags().StringVar(&overwrites.OverwriteSecret, "secret", "", "use this secret instead of the secret in the config file")
	rootCmd.PersistentFlags().StringVar(&overwrites.OverwriteSecretName, "secret-name", "", "use this secret name instead of the secret name in the config file")
	rootCmd.PersistentFlags().StringVar(&overwrites.OverwriteSecretEnv, "secret-from-env", "", "use this environment variable instead of the environment var in the config file")
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

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {

	}
}

func initProjectConfig() {
	projectCfg, projectCfgError = config_parser.ParseRepository(projectCfgFile)
	if projectCfgError == nil {
		projectCfg.MergeWithCliArgs(overwrites)
	}
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
