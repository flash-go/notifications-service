package dto

import (
	"github.com/flash-go/sdk/errors"
)

var (
	ErrInvalidFromEmail = errors.New(errors.ErrBadRequest, "invalid_from_email")
	ErrInvalidFromName  = errors.New(errors.ErrBadRequest, "invalid_from_name")
	ErrInvalidSubject   = errors.New(errors.ErrBadRequest, "invalid_subject")
	ErrInvalidToEmail   = errors.New(errors.ErrBadRequest, "invalid_to_email")
	ErrInvalidHtml      = errors.New(errors.ErrBadRequest, "invalid_html")
	ErrInvalidText      = errors.New(errors.ErrBadRequest, "invalid_text")
)
