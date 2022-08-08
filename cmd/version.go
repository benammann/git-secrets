package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version",
	Run: func(cmd *cobra.Command, args []string) {

		isShort, _ := cmd.Flags().GetBool(FlagShort)

		if isShort == false {
			color.Blue(`
  ________.__  __      _________                            __          
 /  _____/|__|/  |_   /   _____/ ____   ___________   _____/  |_  ______
/   \  ___|  \   __\  \_____  \_/ __ \_/ ___\_  __ \_/ __ \   __\/  ___/
\    \_\  \  ||  |    /        \  ___/\  \___|  | \/\  ___/|  |  \___ \ 
 \______  /__||__|   /_______  /\___  >\___  >__|    \___  >__| /____  >
        \/                   \/     \/     \/            \/          \/`)
			fmt.Println("")
		}
		fmt.Println("Version			", version)
		fmt.Println("Commit			", commit)
		fmt.Println("Date			", date)

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolP(FlagShort, "s", false, "Print version info in short format")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
