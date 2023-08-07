package core

import (
	"fmt"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"auth/services/example"
	"auth/services/user"
	"auth/services/app"
)

func (c *Core) NewDatabase() *gorm.DB {
	c.Log.Info().Msg("Setup database")

	// convert string to integer
	portInt, err := strconv.Atoi(c.Conf.Database.Port)
	if err != nil {
		c.Log.Error().Err(err).Msg("String to int error")
	}
	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Conf.Database.Host, portInt, c.Conf.Database.Username, c.Conf.Database.Password, c.Conf.Database.Name)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		c.Log.Fatal().Err(err).Msg("Error opening db")
	} else {
		c.Log.Info().Msg("DB OK")
	}

	// micgrate cli command is not needed with automigrate here
	db.AutoMigrate(
		example.Example{},
		user.User{},
		app.App{},
	)

	c.registerShutdownFunc(func() error {
		c.Log.Debug().Msg("Closing database connection")
		// Close DB Connection
		sqlDB, err := db.DB()
		sqlDB.Close()

		if err != nil {
			//logger.Errorf("failed to disconnect from database: %v", err)
			return err
		}

		c.Log.Debug().Msg("Database connection closed")
		return nil
	})

	return db

}
