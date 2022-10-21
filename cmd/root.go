package cmd

import (
	"context"
	config_const "github.com/benammann/git-secrets/pkg/config/const"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	global_config "github.com/benammann/git-secrets/pkg/config/global"
	"github.com/benammann/git-secrets/pkg/render"
	"github.com/benammann/git-secrets/pkg/utility"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
	"time"
)

var fs = afero.NewOsFs()

var version string
var commit string
var date string

var globalCfgFile string
var globalCfg *global_config.GlobalConfigProvider
var projectCfgFile string
var projectCfg *config_generic.Repository
var projectCfgError error
var renderingEngine *render.RenderingEngine

var selectedContext *config_generic.Context
var contextName string

var overwrittenSecrets []string

var waitGroupSize = 0

const FlagValue = "value"
const FlagResourceId = "resourceId"
const FlagForce = "force"
const FlagDebug = "debug"
const FlagDryRun = "dry-run"
const FlagTarget = "target"
const FlagAll = "all"
const FlagVerbose = "verbose"
const FlagShort = "short"

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
func Execute(buildVersion string, buildCommit string, buildDate string) {
	version = buildVersion
	commit = buildCommit
	date = buildDate

	ctx := utility.NewChannelContext()
	addChannel, doneChannel := utility.GetContextChannels(ctx)
	ctx, cancel := context.WithCancel(ctx)

	var waitGroup sync.WaitGroup
	timeout, _ := context.WithTimeout(context.Background(), time.Second * 5)

	go func() {
		for {
			select {
			case val := <-addChannel:
				waitGroupSize = waitGroupSize + val
				waitGroup.Add(val)
			case <-doneChannel:
				waitGroup.Done()
				waitGroupSize--
			case <-timeout.Done():
				for waitGroupSize > 0 {
					waitGroup.Done()
					waitGroupSize--
				}
			}
		}
	}()

	cobra.CheckErr(rootCmd.ExecuteContext(ctx))
	cancel()

	waitGroup.Wait()

}

func init() {
	cobra.OnInitialize(initGlobalConfig, initProjectConfig, resolveContext, createRenderingEngine)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&globalCfgFile, "global-config", "", "Path to the global config file: ~/.git-secrets.yaml")
	rootCmd.PersistentFlags().StringVarP(&projectCfgFile, "config", "f", ".git-secrets.json", "Path to the projects config file: .git-secrets.json")
	rootCmd.PersistentFlags().StringVarP(&contextName, "context", "c", "", "Which context to use: default")
	rootCmd.PersistentFlags().StringArrayVar(&overwrittenSecrets, "secret", []string{}, "Pass global secrets directly: --secret secretKey=secretValue")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

// initGlobalConfig reads in config file and ENV variables if set.
func initGlobalConfig() {

	customViper := viper.New()
	customViper.SetFs(afero.NewOsFs())

	shouldCheckErr := len(overwrittenSecrets) == 0

	if globalCfgFile != "" {
		// Use config file from the flag.
		customViper.SetConfigFile(globalCfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if shouldCheckErr == true {
			cobra.CheckErr(err)
		} else if err != nil {
			return
		}

		// Search config in home directory with name ".git-secrets" (without extension).
		customViper.AddConfigPath(home)
		customViper.SetConfigType("yaml")
		customViper.SetConfigName(".git-secrets")
	}

	customViper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if errRead := customViper.ReadInConfig(); errRead == nil {
		cobra.CheckErr(errRead)
	}
	if _, err := os.Stat(customViper.ConfigFileUsed()); err == nil {
		customViper.WatchConfig()
	}

	viperConfigStorage := global_config.NewViperConfigStorage(customViper)
	globalCfg = global_config.NewGlobalConfigProvider(viperConfigStorage)

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

	projectCfg, projectCfgError = config_generic.ParseRepository(fs, projectCfgFile, globalCfg, overwrittenSecretsMap)

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
		renderingEngine = render.NewRenderingEngine(projectCfg, fs, fs)
	}
}
