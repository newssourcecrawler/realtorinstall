package repos

import "errors"

// ErrNotFound is returned by any repo method when a requested record does not exist.
var ErrNotFound = errors.New("repo: not found")
var ErrNameEmailNotFound = errors.New("name and email are required")
var ErrIDNotFound = errors.New("invalid property ID")
var ErrAddrNotFound = errors.New("address and city cannot be empty")
var ErrInvalidCredentials = errors.New("invalid username or password")
var ErrInvalidRegistration = errors.New("username, role, first name, and last name are required")
var ErrUserAlreadyExists = errors.New("username already taken")
var ErrGenerateFromPassword = errors.New("password generation error")
var ErrSignedString = errors.New("signed string error")
var ErrParseWithClaims = errors.New("sparse with claims error")
var ErrInvalidTokenClaims = errors.New("invalid token claims")
var ErrCreateInstallmentPlanIDReq = errors.New("plan_id is required")
