package port

import (
	"github.com/flash-go/sdk/errors"
)

var (
	// Folders
	ErrFolderInvalidParent      = errors.New(errors.ErrBadRequest, "invalid_parent")
	ErrFolderInvalidName        = errors.New(errors.ErrBadRequest, "invalid_name")
	ErrFolderInvalidDescription = errors.New(errors.ErrBadRequest, "invalid_description")
	// Emails
	ErrEmailInvalidId          = errors.New(errors.ErrBadRequest, "invalid_id")
	ErrEmailInvalidFolderId    = errors.New(errors.ErrBadRequest, "invalid_folder_id")
	ErrEmailInvalidFromEmail   = errors.New(errors.ErrBadRequest, "invalid_from_email")
	ErrEmailInvalidFromName    = errors.New(errors.ErrBadRequest, "invalid_from_name")
	ErrEmailInvalidSubject     = errors.New(errors.ErrBadRequest, "invalid_subject")
	ErrEmailInvalidToEmail     = errors.New(errors.ErrBadRequest, "invalid_to_email")
	ErrEmailInvalidHtml        = errors.New(errors.ErrBadRequest, "invalid_html")
	ErrEmailInvalidText        = errors.New(errors.ErrBadRequest, "invalid_text")
	ErrEmailInvalidDescription = errors.New(errors.ErrBadRequest, "invalid_description")
)
