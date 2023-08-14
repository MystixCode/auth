package core

import (
	"auth/services"
	"auth/services/example"
	"auth/services/health"
	"auth/services/root"
	"auth/services/user"
	"auth/services/app"
)

func (c *Core) newServices() *services.Services {

	c.Log.Info().Msg("Setup services")
	rootService := root.NewService()
	healthService := health.NewService(&c.state)
	userService := user.NewService(c.Log, c.Conf, c.Database, c.Validator)
	exampleService := example.NewService(c.Log, c.Conf, c.Database)
	appService := app.NewService(c.Log, c.Conf, c.Database)

	return &services.Services{
		Root:		rootService,
		Health:		healthService,
		User:		userService,
		Example:	exampleService,
		App:		appService,
	}
}
