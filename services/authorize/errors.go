package authorize

import "errors"

var (
	ErrIdParseFailed            = errors.New("failed to parse AuthorizationCode id")
	ErrInsertFailed             = errors.New("failed to insert AuthorizationCode to database")
	ErrNotFound                 = errors.New("AuthorizationCode(s) not found")
	ErrFindFailed               = errors.New("internal error while find AuthorizationCode(s)")
	ErrUpdatedFailed            = errors.New("error while updating AuthorizationCode")
	ErrDeleteFailed             = errors.New("error while deleting AuthorizationCode")
	// ErrPasswordChangeNotAllowed = errors.New("password update only allowed on the 'password_change' endpoint")
)
