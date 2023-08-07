package example

import "errors"

var (
	ErrIdParseFailed = errors.New("failed to parse example id")
	ErrInsertFailed  = errors.New("failed to insert example to database")
	ErrNotFound      = errors.New("example(s) not found")
	ErrFindFailed    = errors.New("internal error while find example(s)")
	ErrUpdatedFailed = errors.New("error while updating example")
	ErrDeleteFailed  = errors.New("error while deleting example")
)
