package port

import (
	"github.com/flash-go/sdk/errors"
)

var (
	// Folders
	ErrFolderNotFound = errors.New(errors.ErrBadRequest, "folder_not_found")
	// Emails
	ErrEmailNotFound = errors.New(errors.ErrBadRequest, "email_not_found")
)
