package conf

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	File     string
	Server   serverConfig
	Log      logConfig
	Database databaseConfig
}

type serverConfig struct {
	Host string
	Port int
}

type logConfig struct {
	Debug bool
	File  string
}

type databaseConfig struct {
	Driver   string
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

func New(cfgFile string) (*Config, error) {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app/")
	viper.AddConfigPath(home + "/.auth/")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if err != nil {
			//fmt.Println("config file not found -> create new file from defaults")

			//create dir if not exists
			os.MkdirAll(filepath.Dir(cfgFile), os.ModePerm)
			//setDefaults
			SetDefaults()
			//Write defaults to config file
			viper.WriteConfig()
		}
	}
	c := Load()

	return c, nil
}

func Load() *Config {
	//setDefaults()

	// 	level, err := log.ParseLevel(viper.GetString("log.level"))
	// 	if err != nil {
	// 		fmt.Println("invalid \"log.level\" config value setting to default \"error\"")
	// 	}

	return &Config{
		Server: serverConfig{
			Host: viper.GetString("server.host"),
			Port: viper.GetInt("server.port"),
		},
		Log: logConfig{
			Debug: viper.GetBool("log.debug"),
			File:  viper.GetString("log.file"),
		},
		Database: databaseConfig{
			Driver:   viper.GetString("database.driver"),
			Username: viper.GetString("database.username"),
			Password: viper.GetString("database.password"),
			Host:     viper.GetString("database.host"),
			Port:     viper.GetString("database.port"),
			Name:     viper.GetString("database.name"),
		},
	}
}
