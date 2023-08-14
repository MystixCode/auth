package user

import "errors"

var (
	ErrIdParseFailed            = errors.New("failed to parse user id")
	ErrInsertFailed             = errors.New("failed to insert user to database")
	ErrNotFound                 = errors.New("user(s) not found")
	ErrFindFailed               = errors.New("internal error while find user(s)")
	ErrUpdatedFailed            = errors.New("error while updating user")
	ErrDeleteFailed             = errors.New("error while deleting user")
	ErrPasswordChangeNotAllowed = errors.New("password update only allowed on the 'password_change' endpoint")
	ErrValidationFailed			= errors.New("validation failed")
)
