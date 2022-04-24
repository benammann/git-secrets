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
	"github.com/spf13/cobra"
)

// addContextCmd represents the addContext command
var addContextCmd = &cobra.Command{
	Use:     "add-context",
	Short:   "adds a context to the existing config file",
	Example: "git-secrets add-context <contextName>",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		writer := projectCfg.GetConfigWriter()
		errAdd := writer.AddContext(args[0])
		cobra.CheckErr(errAdd)
		fmt.Printf("the context %s has been added to your config file\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(addContextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addContextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addContextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
