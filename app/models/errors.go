package models

import "github.com/pkg/errors"

var (
	ErrUserNotPermission = errors.New("no permission")
	ErrUserInvalid       = errors.New("invalid user or not found")
	ErrBadData           = errors.New("bad data")
	ErrTenderNotFound    = errors.New("Tender not found")
	ErrBidNotFound       = errors.New("Bid not found")
)
