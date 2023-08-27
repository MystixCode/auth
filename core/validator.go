package core

import (
	// "os"
	"github.com/go-playground/validator/v10"
	// "git.bitcubix.io/go/validation"
)

// use a single instance of Validate, it caches struct info
//var validate *validator.Validate

// func (c *Core) newValidator() *validation.Validator {
func (c *Core) NewValidator() *validator.Validate {

	c.Log.Info().Msg("Setup validator")
	
	validator := validator.New()

	// Todo: use validator in example service

	return validator
}
