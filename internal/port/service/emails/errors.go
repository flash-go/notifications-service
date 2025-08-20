package port

import (
	"github.com/flash-go/sdk/errors"
)

var (
	// Folders
	ErrFolderExist = errors.New(errors.ErrBadRequest, "folder_exist")
	// Emails
	ErrEmailNotFound = errors.New(errors.ErrBadRequest, "email_not_found")
)
