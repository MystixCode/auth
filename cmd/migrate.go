package cmd

import (
	"fmt"
	"os"

	"auth/core"
	"auth/services/example"
	"auth/services/user"
	"auth/services/app"
	"auth/services/key"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database",
	Long:  `Creates database structure based on models`,
	Run:   migrate,
}

func migrate(cmd *cobra.Command, args []string) {

	appCore, err := core.New(cfgFile, isDebug, logFile)
	if err != nil {
		fmt.Println("core.New() error")
		os.Exit(2)
	}

	appCore.NewDatabase()
	appCore.Database.AutoMigrate(
		example.Example{},
		user.User{},
		app.App{},
		key.Key{},
	)

}
