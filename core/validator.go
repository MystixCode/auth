package core

import (
	// "os"
	// "github.com/go-playground/validator/v10"
	// "git.bitcubix.io/go/validation"
)

// func (c *Core) newValidator() *validation.Validator {
	func (c *Core) newValidator() {

	// validator := validation.NewValidator()

	// Todo: add return err to validation.NewValidator()
	// validator, err := validation.NewValidator()
	// if err != nil {
	// 	log.Error().Err(err).Msg("Setup validator error")
	// 	os.Exit(2)
	// }

	// Todo: Setup Validator --> check stuff in git.bitcubix.io/go/validation
	// Todo: use validator in example and user service

	c.Log.Info().Msg("Setup validator")
	c.Log.Warn().Msg("TODO: Setup validator")
	// return validator
	return
}
