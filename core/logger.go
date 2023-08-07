package core

import (
	"os"

	"auth/log"
)

func (c *Core) newLog(isDebug bool, logFile string) *log.Logger {

	//fmt.Println("Input logfile: ", logFile)
	//fmt.Println("Input isDebug: ", isDebug)

	//Set TimeZone
	os.Setenv("TZ", "Europe/Zurich")

	if logFile == "" {
		logFile = c.Conf.Log.File
	} else {
		c.Conf.Log.File = logFile
	}

	if !isDebug {
		//fmt.Println("debug not set by cli -> get it from config")
		isDebug = c.Conf.Log.Debug
	}

	log, err := log.New(isDebug, logFile)
	if err != nil {
		log.Error().Err(err).Msg("Setup log error")
		os.Exit(2)
	}
	log.Info().Msg("Setup conf")
	log.Debug().Str("Host", c.Conf.Server.Host).Int("Port", c.Conf.Server.Port).Msg("Server")
	log.Debug().Str("File", c.Conf.Log.File).Bool("Debug", c.Conf.Log.Debug).Msg("Logger")
	log.Debug().Str("Host", c.Conf.Database.Host).Str("Username", c.Conf.Database.Username).Str("Password", c.Conf.Database.Password).Str("Driver", c.Conf.Database.Driver).Str("Port", c.Conf.Database.Port).Str("Name", c.Conf.Database.Name).Msg("Database")
	log.Info().Msg("Setup log")
	return log
}
