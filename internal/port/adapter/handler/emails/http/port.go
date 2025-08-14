package port

import (
	"github.com/flash-go/flash/http/server"
)

type Interface interface {
	Send(ctx server.ReqCtx)
}
