package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of auth",
	Long:  `All software has versions. This is auth's`,
	Run:   version,
}

func version(cmd *cobra.Command, args []string) {
	fmt.Println("auth version 0.0.1")
}
