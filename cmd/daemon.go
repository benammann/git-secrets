/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	cli_config "github.com/benammann/git-secrets/pkg/config/cli"
	"github.com/benammann/git-secrets/pkg/config/daemon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Args: cobra.MaximumNArgs(2),
	Example: `
git-secrets daemon
git-secrets daemon <file-to-watch> <context>
`,
	Short: "automatically render files on config changes",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 2 {
			 file := args[0]
			 context := args[1]
			 absPath, errAbs := filepath.Abs(file)
			 cobra.CheckErr(errAbs)
			 viper.Set(fmt.Sprintf("%s.%s", cli_config.DaemonWatches, absPath), context)
			 fmt.Println(absPath, "added to watches")
		} else {
			daemonInstance := daemon.NewDaemon()
			fmt.Println(daemonInstance.Run())
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
