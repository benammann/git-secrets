/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	gcp_auth "github.com/benammann/git-secrets/pkg/gcp/auth"
	"github.com/spf13/cobra"
)


var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands for remote secret managers",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var authGcpCmd = &cobra.Command{
	Use: "gcp",
	Short: "Authenticate against GoogleCloud using gcloud-cli",
	RunE: func(cmd *cobra.Command, args []string) error {

		force, _ := cmd.Flags().GetBool(FlagForce)

		isAuthenticated, errAuth := gcp_auth.IsAuthenticated()
		cobra.CheckErr(errAuth)

		if isAuthenticated && force == false {
			fmt.Println("you are already authenticated. Use --force if you still want to continue")
			return nil
		}

		cobra.CheckErr(gcp_auth.Authenticate())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.AddCommand(authGcpCmd)
	authGcpCmd.Flags().Bool(FlagForce, false, "Use --force to ignore existing gcp authentication")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
