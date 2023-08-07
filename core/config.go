package core

import (
	"auth/conf"

	"fmt"
	"os"
)

func (c *Core) newConf(cfgFile string) *conf.Config {
	conf, err := conf.New(cfgFile)
	if err != nil {
		fmt.Println("Setup config error")
		os.Exit(2)
	}
	return conf

}
