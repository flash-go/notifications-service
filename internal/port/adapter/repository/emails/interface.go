package port

import (
	"context"
)

type Interface interface {
	// Folders
	CreateFolder(ctx context.Context, data CreateFolderData) (*FolderResult, error)
	FilterFolders(ctx context.Context, data FilterFoldersData) (*[]FolderResult, error)
	DeleteFolder(ctx context.Context, id uint) error
	UpdateFolder(ctx context.Context, id uint, data map[string]any) error
	// Emails
	CreateEmail(ctx context.Context, data CreateEmailData) (*EmailResult, error)
	FilterEmails(ctx context.Context, data FilterEmailsData) (*[]EmailResult, error)
	DeleteEmail(ctx context.Context, id uint) error
	UpdateEmail(ctx context.Context, id uint, data map[string]any) error
	Send(ctx context.Context, data SendData) (*EmailLogResult, error)
	FilterEmailLogs(ctx context.Context, data FilterEmailLogsData) (*[]EmailLogResult, error)
}
