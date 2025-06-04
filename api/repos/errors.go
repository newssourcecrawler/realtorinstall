package repos

import "errors"

// ErrNotFound is returned by any repo method when a requested record does not exist.
var ErrNotFound = errors.New("repo: not found")
var NameEmailNotFound = errors.New("name and email are required")
var IDNotFound = errors.New("invalid property ID")
var AddrNotFound = errors.New("address and city cannot be empty")
