package conf

import "github.com/spf13/viper"

func SetDefaults() {
	// server
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8080")

	// database
	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.username", "dbuser")
	viper.SetDefault("database.password", "changeit")
	viper.SetDefault("database.host", "pg")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.name", "auth")

	// logger
	viper.SetDefault("log.debug", false)
	viper.SetDefault("log.file", "/var/log/auth/auth.log")
}
