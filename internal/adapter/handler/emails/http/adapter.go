package adapter

import (
	"github.com/flash-go/flash/http/server"
	dto "github.com/flash-go/notifications-service/internal/dto/emails"
	httpEmailsHandlerAdapterPort "github.com/flash-go/notifications-service/internal/port/adapter/handler/emails/http"
	emailsServicePort "github.com/flash-go/notifications-service/internal/port/service/emails"
	"github.com/flash-go/sdk/errors"
)

type Config struct {
	EmailsService emailsServicePort.Interface
}

func New(config *Config) httpEmailsHandlerAdapterPort.Interface {
	return &adapter{
		config.EmailsService,
	}
}

type adapter struct {
	emailsService emailsServicePort.Interface
}

// @Summary Send email
// @Tags emails
// @Accept json
// @Produce json,plain
// @Param request body dto.SendRequest true "Send email"
// @Success 201 {object} dto.EmailResponse
// @Failure 400 {string} string "Possible error codes: bad_request, bad_request:invalid_path"
// @Router /emails/send [post]
func (a *adapter) Send(ctx server.ReqCtx) {
	// Parse request json body
	var request dto.SendRequest
	if err := ctx.ReadJson(&request); err != nil {
		ctx.WriteErrorResponse(errors.ErrBadRequest)
		return
	}

	// Validate request
	if err := request.Validate(); err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Create data
	data := emailsServicePort.SendData(request)

	// Send email
	email, err := a.emailsService.Send(
		ctx.Context(),
		&data,
	)
	if err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Write success response
	ctx.WriteResponse(201, dto.EmailResponse(*email))
}
