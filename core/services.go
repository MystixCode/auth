package core

import (
	"auth/services"
	"auth/services/example"
	"auth/services/health"
	"auth/services/root"
	"auth/services/user"
	"auth/services/app"
	"auth/services/key"
)

func (c *Core) newServices() *services.Services {

	c.Log.Info().Msg("Setup services")
	rootService := root.NewService()
	healthService := health.NewService(&c.state)
	userService := user.NewService(c.Log, c.Conf, c.Database, c.Validator)
	exampleService := example.NewService(c.Log, c.Conf, c.Database)
    // Create an instance of key.Service
    keyService := key.NewService(c.Log, c.Conf, c.Database)

    // Create an instance of app.Service and inject the key.Service instance
    appService := app.NewService(c.Log, c.Conf, c.Database, keyService)

	return &services.Services{
		Root:		rootService,
		Health:		healthService,
		User:		userService,
		Example:	exampleService,
		App:		appService,
		Key:		keyService,
	}
}
