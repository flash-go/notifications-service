package port

import (
	"github.com/flash-go/flash/http/server"
)

type Interface interface {
	// Folders
	AdminCreateFolder(ctx server.ReqCtx)
	AdminFilterFolders(ctx server.ReqCtx)
	AdminDeleteFolder(ctx server.ReqCtx)
	AdminUpdateFolder(ctx server.ReqCtx)
	// Emails
	AdminCreateEmail(ctx server.ReqCtx)
	AdminFilterEmails(ctx server.ReqCtx)
	AdminDeleteEmail(ctx server.ReqCtx)
	AdminUpdateEmail(ctx server.ReqCtx)
	SendCustom(ctx server.ReqCtx)
	Send(ctx server.ReqCtx)
	AdminFilterEmailLogs(ctx server.ReqCtx)
}
