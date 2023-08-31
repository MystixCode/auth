package auth

import "errors"

var (
	ErrIdParseFailed    = errors.New("failed to parse app id")
	ErrInsertFailed     = errors.New("failed to insert app to database")
	ErrNotFound         = errors.New("app(s) not found")
	ErrFindFailed       = errors.New("internal error while find app(s)")
	ErrUpdatedFailed    = errors.New("error while updating app")
	ErrDeleteFailed		= errors.New("error while deleting app")
	ErrValidationFailed	= errors.New("validation failed")
	// ErrPasswordChangeNotAllowed = errors.New("password update only allowed on the 'password_change' endpoint")

	// ErrIdParseFailed 	= errors.New("failed to parse key id")
	// ErrInsertFailed  	= errors.New("failed to insert key to database")
	// ErrNotFound      	= errors.New("key(s) not found")
	// ErrFindFailed    	= errors.New("internal error while find key(s)")
	// ErrDeleteFailed  	= errors.New("error while deleting key")
	ErrKeyGenFailed		= errors.New("error while generating key")
	// ErrValidationFailed	= errors.New("validation failed")


	// ErrIdParseFailed            = errors.New("failed to parse user id")
	// ErrInsertFailed             = errors.New("failed to insert user to database")
	// ErrNotFound                 = errors.New("user(s) not found")
	// ErrFindFailed               = errors.New("internal error while find user(s)")
	// ErrUpdatedFailed            = errors.New("error while updating user")
	// ErrDeleteFailed             = errors.New("error while deleting user")
	ErrPasswordChangeNotAllowed = errors.New("password update only allowed on the 'password_change' endpoint")
	// ErrValidationFailed			= errors.New("validation failed")
	ErrLogin					= errors.New("login failed")

)