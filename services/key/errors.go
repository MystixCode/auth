package key

import "errors"

var (
	ErrIdParseFailed = errors.New("failed to parse key id")
	ErrInsertFailed  = errors.New("failed to insert key to database")
	ErrNotFound      = errors.New("key(s) not found")
	ErrFindFailed    = errors.New("internal error while find key(s)")
	ErrUpdatedFailed = errors.New("error while updating key")
	ErrDeleteFailed  = errors.New("error while deleting key")
)
