package core

import "auth/api"

func (c *Core) newApi() *api.Api {

	c.Log.Info().Msg("Setup api")

	api := &api.Api{
		Root:   	api.NewRootEndpoint(c.Log, c.services.Root),
		Health:  	api.NewHealthEndpoint(c.Log, c.services.Health),
		User:    	api.NewUserEndpoint(c.Log, c.services.User),
		Example:	api.NewExampleEndpoint(c.Log, c.services.Example),
		App: 		api.NewAppEndpoint(c.Log, c.services.App),
	}

	api.New(c.router)

	return api
}
