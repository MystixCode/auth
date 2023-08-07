package app

import "errors"

var (
	ErrIdParseFailed            = errors.New("failed to parse app id")
	ErrInsertFailed             = errors.New("failed to insert app to database")
	ErrNotFound                 = errors.New("app(s) not found")
	ErrFindFailed               = errors.New("internal error while find app(s)")
	ErrUpdatedFailed            = errors.New("error while updating app")
	ErrDeleteFailed             = errors.New("error while deleting app")
	// ErrPasswordChangeNotAllowed = errors.New("password update only allowed on the 'password_change' endpoint")
)
