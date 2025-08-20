package adapter

import (
	"github.com/flash-go/flash/http/server"
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

// Folders

// @Summary Create email folder (admin)
// @Tags emails
// @Security BearerAuth
// @Accept json
// @Produce json,plain
// @Param request body httpEmailsHandlerAdapterPort.CreateFolderData true "Create email folder"
// @Success 201 {object} httpEmailsHandlerAdapterPort.FolderResponse
// @Failure 400 {string} string "Possible error codes: bad_request, bad_request:invalid_parent, bad_request:invalid_name, bad_request:folder_exist"
// @Router /admin/notifications/emails/folders [post]
func (a *adapter) AdminCreateFolder(ctx server.ReqCtx) {
	// Create email folder
	folder, err := a.emailsService.CreateFolder(
		ctx.Context(),
		emailsServicePort.CreateFolderData(
			*ctx.GetJsonBody().(*httpEmailsHandlerAdapterPort.CreateFolderData),
		),
	)
	if err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Write success response
	ctx.WriteResponse(201, httpEmailsHandlerAdapterPort.FolderResponse(*folder))
}

// @Summary Filter email folders (admin)
// @Tags emails
// @Security BearerAuth
// @Accept json
// @Produce json,plain
// @Param request body httpEmailsHandlerAdapterPort.FilterFoldersData true "Filter email folders"
// @Success 200 {array} httpEmailsHandlerAdapterPort.FolderResponse
// @Failure 400 {string} string "Possible error codes: bad_request"
// @Router /admin/notifications/emails/folders/filter [post]
func (a *adapter) AdminFilterFolders(ctx server.ReqCtx) {
	// Filter email folders
	folders, err := a.emailsService.FilterFolders(
		ctx.Context(),
		emailsServicePort.FilterFoldersData(
			*ctx.GetJsonBody().(*httpEmailsHandlerAdapterPort.FilterFoldersData),
		),
	)
	if err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Map service to adapter results
	results := make([]httpEmailsHandlerAdapterPort.FolderResponse, 0, len(*folders))
	for _, folder := range *folders {
		results = append(
			results,
			httpEmailsHandlerAdapterPort.FolderResponse(folder),
		)
	}

	// Write success response
	ctx.WriteResponse(200, results)
}

// @Summary Delete email folder (admin)
// @Tags emails
// @Security BearerAuth
// @Produce plain
// @Param id path int true "Folder ID"
// @Success 204
// @Failure 400 {string} string "Possible error codes: bad_request, bad_request:folder_not_found"
// @Router /admin/notifications/emails/folders/{id} [delete]
func (a *adapter) AdminDeleteFolder(ctx server.ReqCtx) {
	// Get and convert folder id to uint64
	id, err := ctx.UserValueUint64("id")
	if err != nil {
		ctx.WriteErrorResponse(errors.ErrBadRequest)
		return
	}

	// Delete email folder
	if err := a.emailsService.DeleteFolder(
		ctx.Context(),
		uint(id),
	); err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Write success response
	ctx.WriteResponse(204, nil)
}

// @Summary Update email folder (admin)
// @Tags emails
// @Security BearerAuth
// @Accept json
// @Produce json,plain
// @Param id path int true "Folder ID"
// @Param request body httpEmailsHandlerAdapterPort._UpdateFolderData true "Update email folder"
// @Success 204
// @Failure 400 {string} string "Possible error codes: bad_request, bad_request:invalid_name, bad_request:folder_not_found"
// @Router /admin/notifications/emails/folders/{id} [patch]
func (a *adapter) AdminUpdateFolder(ctx server.ReqCtx) {
	// Get and convert folder id to uint64
	id, err := ctx.UserValueUint64("id")
	if err != nil {
		ctx.WriteErrorResponse(errors.ErrBadRequest)
		return
	}

	// Get data
	data := ctx.GetJsonBody().(*httpEmailsHandlerAdapterPort.UpdateFolderData)

	// Set folder data
	folder := make(map[string]any)

	if data.ParentId.Set {
		folder["parent_id"] = data.ParentId.Value
	}
	if data.Name.Set {
		folder["name"] = data.Name.Value
	}
	if data.Description.Set {
		folder["Description"] = data.Description.Value
	}

	// Update email folder
	if err := a.emailsService.UpdateFolder(
		ctx.Context(),
		uint(id),
		folder,
	); err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Write success response
	ctx.WriteResponse(204, nil)
}

// Emails

// @Summary Create email (admin)
// @Tags emails
// @Security BearerAuth
// @Accept json
// @Produce json,plain
// @Param request body httpEmailsHandlerAdapterPort.CreateEmailData true "Create email"
// @Success 201 {object} httpEmailsHandlerAdapterPort.EmailResponse
// @Failure 400 {string} string "Possible error codes: bad_request, bad_request:invalid_folder_id, bad_request:invalid_from_email, bad_request:invalid_from_name, bad_request:invalid_subject, bad_request:invalid_html, bad_request:invalid_text"
// @Router /admin/notifications/emails [post]
func (a *adapter) AdminCreateEmail(ctx server.ReqCtx) {
	// Create email
	email, err := a.emailsService.CreateEmail(
		ctx.Context(),
		emailsServicePort.CreateEmailData(
			*ctx.GetJsonBody().(*httpEmailsHandlerAdapterPort.CreateEmailData),
		),
	)
	if err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Write success response
	ctx.WriteResponse(201, httpEmailsHandlerAdapterPort.EmailResponse(*email))
}

// @Summary Filter emails (admin)
// @Tags emails
// @Security BearerAuth
// @Accept json
// @Produce json,plain
// @Param request body httpEmailsHandlerAdapterPort.FilterEmailsData true "Filter emails"
// @Success 200 {array} httpEmailsHandlerAdapterPort.EmailResponse
// @Failure 400 {string} string "Possible error codes: bad_request"
// @Router /admin/notifications/emails/filter [post]
func (a *adapter) AdminFilterEmails(ctx server.ReqCtx) {
	// Filter emails
	emails, err := a.emailsService.FilterEmails(
		ctx.Context(),
		emailsServicePort.FilterEmailsData(
			*ctx.GetJsonBody().(*httpEmailsHandlerAdapterPort.FilterEmailsData),
		),
	)
	if err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Map service to adapter results
	results := make([]httpEmailsHandlerAdapterPort.EmailResponse, 0, len(*emails))
	for _, email := range *emails {
		results = append(
			results,
			httpEmailsHandlerAdapterPort.EmailResponse(email),
		)
	}

	// Write success response
	ctx.WriteResponse(200, results)
}

// @Summary Delete email (admin)
// @Tags emails
// @Security BearerAuth
// @Produce plain
// @Param id path int true "Email ID"
// @Success 204
// @Failure 400 {string} string "Possible error codes: bad_request, bad_request:email_not_found"
// @Router /admin/notifications/emails/{id} [delete]
func (a *adapter) AdminDeleteEmail(ctx server.ReqCtx) {
	// Get and convert email id to uint64
	id, err := ctx.UserValueUint64("id")
	if err != nil {
		ctx.WriteErrorResponse(errors.ErrBadRequest)
		return
	}

	// Delete email
	if err := a.emailsService.DeleteEmail(
		ctx.Context(),
		uint(id),
	); err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Write success response
	ctx.WriteResponse(204, nil)
}

// @Summary Update email (admin)
// @Tags emails
// @Security BearerAuth
// @Accept json
// @Produce json,plain
// @Param id path int true "Email ID"
// @Param request body httpEmailsHandlerAdapterPort._UpdateEmailData true "Update email"
// @Success 204
// @Failure 400 {string} string "Possible error codes: bad_request, bad_request:invalid_folder_id, bad_request:invalid_from_email, bad_request:invalid_from_name, bad_request:invalid_subject, bad_request:invalid_html, bad_request:invalid_text, bad_request:email_not_found"
// @Router /admin/notifications/emails/{id} [patch]
func (a *adapter) AdminUpdateEmail(ctx server.ReqCtx) {
	// Get and convert email id to uint64
	id, err := ctx.UserValueUint64("id")
	if err != nil {
		ctx.WriteErrorResponse(errors.ErrBadRequest)
		return
	}

	// Get data
	data := ctx.GetJsonBody().(*httpEmailsHandlerAdapterPort.UpdateEmailData)

	// Set email data
	email := make(map[string]any)

	if data.FolderId.Set {
		email["folder_id"] = data.FolderId.Value
	}
	if data.FromEmail.Set {
		email["from_email"] = data.FromEmail.Value
	}
	if data.FromName.Set {
		email["from_name"] = data.FromName.Value
	}
	if data.Subject.Set {
		email["subject"] = data.Subject.Value
	}
	if data.Html.Set {
		email["html"] = data.Html.Value
	}
	if data.Text.Set {
		email["text"] = data.Text.Value
	}
	if data.Description.Set {
		email["description"] = data.Description.Value
	}

	// Update email
	if err := a.emailsService.UpdateEmail(
		ctx.Context(),
		uint(id),
		email,
	); err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Write success response
	ctx.WriteResponse(204, nil)
}

// @Summary Send custom email
// @Tags emails
// @Security BearerAuth
// @Accept json
// @Produce json,plain
// @Param request body httpEmailsHandlerAdapterPort.SendCustomData true "Send custom email"
// @Success 201 {object} httpEmailsHandlerAdapterPort.EmailLogResponse
// @Failure 400 {string} string "Possible error codes: bad_request, bad_request:invalid_from_email, bad_request:invalid_from_name, bad_request:invalid_subject, bad_request:invalid_to_email, bad_request:invalid_html, bad_request:invalid_text"
// @Router /notifications/emails/send/custom [post]
func (a *adapter) SendCustom(ctx server.ReqCtx) {
	// Send custom email
	log, err := a.emailsService.SendCustom(
		ctx.Context(),
		emailsServicePort.SendCustomData(
			*ctx.GetJsonBody().(*httpEmailsHandlerAdapterPort.SendCustomData),
		),
	)
	if err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Write success response
	ctx.WriteResponse(201, httpEmailsHandlerAdapterPort.EmailLogResponse(*log))
}

// @Summary Send email
// @Tags emails
// @Accept json
// @Produce json,plain
// @Param request body httpEmailsHandlerAdapterPort.SendData true "Send email"
// @Success 201 {object} httpEmailsHandlerAdapterPort.EmailLogResponse
// @Failure 400 {string} string "Possible error codes: bad_request, bad_request:invalid_id, bad_request:invalid_to_email, bad_request:email_not_found"
// @Router /notifications/emails/send [post]
func (a *adapter) Send(ctx server.ReqCtx) {
	// Send email
	log, err := a.emailsService.Send(
		ctx.Context(),
		emailsServicePort.SendData(
			*ctx.GetJsonBody().(*httpEmailsHandlerAdapterPort.SendData),
		),
	)
	if err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Write success response
	ctx.WriteResponse(201, httpEmailsHandlerAdapterPort.EmailLogResponse(*log))
}

// @Summary Filter email logs (admin)
// @Tags emails
// @Security BearerAuth
// @Accept json
// @Produce json,plain
// @Param request body httpEmailsHandlerAdapterPort.FilterEmailLogsData true "Filter email logs (admin)"
// @Success 200 {array} httpEmailsHandlerAdapterPort.EmailLogResponse
// @Failure 400 {string} string "Possible error codes: bad_request"
// @Router /admin/notifications/emails/logs/filter [post]
func (a *adapter) AdminFilterEmailLogs(ctx server.ReqCtx) {
	// Filter email logs
	logs, err := a.emailsService.FilterEmailLogs(
		ctx.Context(),
		emailsServicePort.FilterEmailLogsData(
			*ctx.GetJsonBody().(*httpEmailsHandlerAdapterPort.FilterEmailLogsData),
		),
	)
	if err != nil {
		ctx.WriteErrorResponse(err)
		return
	}

	// Map service to adapter results
	results := make([]httpEmailsHandlerAdapterPort.EmailLogResponse, 0, len(*logs))
	for _, log := range *logs {
		results = append(
			results,
			httpEmailsHandlerAdapterPort.EmailLogResponse(log),
		)
	}

	// Write success response
	ctx.WriteResponse(200, results)
}
