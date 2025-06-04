package repos

import "errors"

// ErrNotFound is returned by any repo method when a requested record does not exist.
var ErrNotFound = errors.New("repo: not found")
