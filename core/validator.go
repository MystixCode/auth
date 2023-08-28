package core

import (

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
//var validate *validator.Validate

func (c *Core) NewValidator() *validator.Validate {

	c.Log.Info().Msg("Setup validator")
	
	validator := validator.New()

	return validator
}
