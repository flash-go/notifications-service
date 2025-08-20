package port

import (
	"encoding/json"
	"net/mail"
	"time"

	"github.com/flash-go/sdk/types"
)

// Data

// Folders

type CreateFolderData struct {
	ParentId    *uint  `json:"parent_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SystemFlag  bool   `json:"system_flag"`
}

func (r *CreateFolderData) Validate() error {
	if err := r.ValidateParentId(); err != nil {
		return err
	}
	if err := r.ValidateName(); err != nil {
		return err
	}
	return nil
}
func (r *CreateFolderData) ValidateParentId() error {
	if r.ParentId != nil && *r.ParentId <= 0 {
		return ErrFolderInvalidParent
	}
	return nil
}
func (r *CreateFolderData) ValidateName() error {
	if r.Name == "" {
		return ErrFolderInvalidName
	}
	return nil
}

type FilterFoldersData struct {
	Id         *[]uint   `json:"id"`
	ParentId   *[]*uint  `json:"parent_id"`
	Name       *[]string `json:"name"`
	SystemFlag *bool     `json:"system_flag"`
}

func (r *FilterFoldersData) Validate() error {
	return nil
}

//lint:ignore U1000 Need for @Param request body in httpEmailsHandlerAdapterPort.AdminUpdateFolder
type _UpdateFolderData struct {
	ParentId    uint   `json:"parent_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateFolderData struct {
	ParentId    types.Nullable[uint]   `json:"parent_id"`
	Name        types.Nullable[string] `json:"name"`
	Description types.Nullable[string] `json:"description"`
}

func (r *UpdateFolderData) Validate() error {
	if err := r.ValidateParentId(); err != nil {
		return err
	}
	if err := r.ValidateName(); err != nil {
		return err
	}
	if err := r.ValidateDescription(); err != nil {
		return err
	}
	return nil
}
func (r *UpdateFolderData) ValidateParentId() error {
	if r.ParentId.Set && r.ParentId.Value != nil && *r.ParentId.Value < 1 {
		return ErrFolderInvalidParent
	}
	return nil
}
func (r *UpdateFolderData) ValidateName() error {
	if r.Name.Set && (r.Name.Value == nil || *r.Name.Value == "") {
		return ErrFolderInvalidName
	}
	return nil
}
func (r *UpdateFolderData) ValidateDescription() error {
	if r.Description.Set && (r.Description.Value == nil || *r.Description.Value == "") {
		return ErrFolderInvalidDescription
	}
	return nil
}

// Emails

type CreateEmailData struct {
	FolderId    *uint  `json:"folder_id"`
	FromEmail   string `json:"from_email"`
	FromName    string `json:"from_name"`
	Subject     string `json:"subject"`
	Html        string `json:"html"`
	Text        string `json:"text"`
	Description string `json:"description"`
	SystemFlag  bool   `json:"system_flag"`
}

func (r *CreateEmailData) Validate() error {
	if err := r.ValidateFolderId(); err != nil {
		return err
	}
	if err := r.ValidateFromEmail(); err != nil {
		return err
	}
	if err := r.ValidateFromName(); err != nil {
		return err
	}
	if err := r.ValidateSubject(); err != nil {
		return err
	}
	if err := r.ValidateHtml(); err != nil {
		return err
	}
	if err := r.ValidateText(); err != nil {
		return err
	}
	return nil
}
func (r *CreateEmailData) ValidateFolderId() error {
	if r.FolderId != nil && *r.FolderId <= 0 {
		return ErrEmailInvalidFolderId
	}
	return nil
}
func (r *CreateEmailData) ValidateFromEmail() error {
	if r.FromEmail == "" {
		return ErrEmailInvalidFromEmail
	}
	if _, err := mail.ParseAddress(r.FromEmail); err != nil {
		return ErrEmailInvalidToEmail
	}
	return nil
}
func (r *CreateEmailData) ValidateFromName() error {
	if r.FromName == "" {
		return ErrEmailInvalidFromName
	}
	return nil
}
func (r *CreateEmailData) ValidateSubject() error {
	if r.Subject == "" {
		return ErrEmailInvalidSubject
	}
	return nil
}
func (r *CreateEmailData) ValidateHtml() error {
	if r.Html == "" {
		return ErrEmailInvalidHtml
	}
	return nil
}
func (r *CreateEmailData) ValidateText() error {
	if r.Text == "" {
		return ErrEmailInvalidText
	}
	return nil
}

type FilterEmailsData struct {
	Id         *[]uint  `json:"id"`
	FolderId   *[]*uint `json:"folder_id"`
	SystemFlag *bool    `json:"system_flag"`
}

func (r *FilterEmailsData) Validate() error {
	return nil
}

//lint:ignore U1000 Need for @Param request body in httpEmailsHandlerAdapterPort.AdminUpdateEmail
type _UpdateEmailData struct {
	FolderId    uint   `json:"folder_id"`
	FromEmail   string `json:"from_email"`
	FromName    string `json:"from_name"`
	Subject     string `json:"subject"`
	Html        string `json:"html"`
	Text        string `json:"text"`
	Description string `json:"description"`
}

type UpdateEmailData struct {
	FolderId    types.Nullable[uint]   `json:"folder_id"`
	FromEmail   types.Nullable[string] `json:"from_email"`
	FromName    types.Nullable[string] `json:"from_name"`
	Subject     types.Nullable[string] `json:"subject"`
	Html        types.Nullable[string] `json:"html"`
	Text        types.Nullable[string] `json:"text"`
	Description types.Nullable[string] `json:"description"`
}

func (r *UpdateEmailData) Validate() error {
	if err := r.ValidateFolderId(); err != nil {
		return err
	}
	if err := r.ValidateFromEmail(); err != nil {
		return err
	}
	if err := r.ValidateFromName(); err != nil {
		return err
	}
	if err := r.ValidateSubject(); err != nil {
		return err
	}
	if err := r.ValidateHtml(); err != nil {
		return err
	}
	if err := r.ValidateText(); err != nil {
		return err
	}
	if err := r.ValidateDescription(); err != nil {
		return err
	}
	return nil
}
func (r *UpdateEmailData) ValidateFolderId() error {
	if r.FolderId.Set && r.FolderId.Value != nil && *r.FolderId.Value < 1 {
		return ErrEmailInvalidFolderId
	}
	return nil
}
func (r *UpdateEmailData) ValidateFromEmail() error {
	if r.FromEmail.Set {
		if r.FromEmail.Value == nil || *r.FromEmail.Value == "" {
			return ErrEmailInvalidFromEmail
		}
		if _, err := mail.ParseAddress(*r.FromEmail.Value); err != nil {
			return ErrEmailInvalidToEmail
		}
	}
	return nil
}
func (r *UpdateEmailData) ValidateFromName() error {
	if r.FromName.Set && (r.FromName.Value == nil || *r.FromName.Value == "") {
		return ErrEmailInvalidFromName
	}
	return nil
}
func (r *UpdateEmailData) ValidateSubject() error {
	if r.Subject.Set && (r.Subject.Value == nil || *r.Subject.Value == "") {
		return ErrEmailInvalidSubject
	}
	return nil
}
func (r *UpdateEmailData) ValidateHtml() error {
	if r.Html.Set && (r.Html.Value == nil || *r.Html.Value == "") {
		return ErrEmailInvalidHtml
	}
	return nil
}
func (r *UpdateEmailData) ValidateText() error {
	if r.Text.Set && (r.Text.Value == nil || *r.Text.Value == "") {
		return ErrEmailInvalidText
	}
	return nil
}
func (r *UpdateEmailData) ValidateDescription() error {
	if r.Description.Set && (r.Description.Value == nil || *r.Description.Value == "") {
		return ErrEmailInvalidDescription
	}
	return nil
}

type SendCustomData struct {
	FromEmail string `json:"from_email"`
	FromName  string `json:"from_name"`
	Subject   string `json:"subject"`
	ToEmail   string `json:"to_email"`
	Html      string `json:"html"`
	Text      string `json:"text"`
}

func (r *SendCustomData) Validate() error {
	if err := r.ValidateFromEmail(); err != nil {
		return err
	}
	if err := r.ValidateFromName(); err != nil {
		return err
	}
	if err := r.ValidateSubject(); err != nil {
		return err
	}
	if err := r.ValidateToEmail(); err != nil {
		return err
	}
	if err := r.ValidateHtml(); err != nil {
		return err
	}
	if err := r.ValidateText(); err != nil {
		return err
	}
	return nil
}
func (r *SendCustomData) ValidateFromEmail() error {
	if r.FromEmail == "" {
		return ErrEmailInvalidFromEmail
	}
	return nil
}
func (r *SendCustomData) ValidateFromName() error {
	if r.FromName == "" {
		return ErrEmailInvalidFromName
	}
	return nil
}
func (r *SendCustomData) ValidateSubject() error {
	if r.Subject == "" {
		return ErrEmailInvalidSubject
	}
	return nil
}
func (r *SendCustomData) ValidateToEmail() error {
	if r.ToEmail == "" {
		return ErrEmailInvalidToEmail
	}
	if _, err := mail.ParseAddress(r.ToEmail); err != nil {
		return ErrEmailInvalidToEmail
	}
	return nil
}
func (r *SendCustomData) ValidateHtml() error {
	if r.Html == "" {
		return ErrEmailInvalidHtml
	}
	return nil
}
func (r *SendCustomData) ValidateText() error {
	if r.Text == "" {
		return ErrEmailInvalidText
	}
	return nil
}

type SendData struct {
	EmailId uint             `json:"email_id"`
	ToEmail string           `json:"to_email"`
	Vars    *json.RawMessage `json:"vars"`
}

func (r *SendData) Validate() error {
	if err := r.ValidateEmailId(); err != nil {
		return err
	}
	if err := r.ValidateToEmail(); err != nil {
		return err
	}
	return nil
}
func (r *SendData) ValidateEmailId() error {
	if r.EmailId <= 0 {
		return ErrEmailInvalidId
	}
	return nil
}
func (r *SendData) ValidateToEmail() error {
	if r.ToEmail == "" {
		return ErrEmailInvalidToEmail
	}
	if _, err := mail.ParseAddress(r.ToEmail); err != nil {
		return ErrEmailInvalidToEmail
	}
	return nil
}

type FilterEmailLogsData struct {
	Id        *[]uint   `json:"id"`
	FromEmail *[]string `json:"from_email"`
	FromName  *[]string `json:"from_name"`
	ToEmail   *[]string `json:"to_email"`
	Status    *[]string `json:"status"`
	MessageId *[]string `json:"message_id"`
}

func (r *FilterEmailLogsData) Validate() error {
	return nil
}

// Responses

type FolderResponse struct {
	Id          uint      `json:"id"`
	ParentId    *uint     `json:"parent_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SystemFlag  bool      `json:"system_flag"`
	Updated     time.Time `json:"updated"`
	Created     time.Time `json:"created"`
}

type EmailResponse struct {
	Id          uint      `json:"id"`
	FolderId    *uint     `json:"folder_id"`
	FromEmail   string    `json:"from_email"`
	FromName    string    `json:"from_name"`
	Subject     string    `json:"subject"`
	Html        string    `json:"html"`
	Text        string    `json:"text"`
	Description string    `json:"description"`
	SystemFlag  bool      `json:"system_flag"`
	Updated     time.Time `json:"updated"`
	Created     time.Time `json:"created"`
}

type EmailLogResponse struct {
	Id        uint      `json:"id"`
	FromEmail string    `json:"from_email"`
	FromName  string    `json:"from_name"`
	Subject   string    `json:"subject"`
	ToEmail   string    `json:"to_email"`
	Html      string    `json:"html"`
	Text      string    `json:"text"`
	Status    string    `json:"status"`
	MessageId *string   `json:"message_id"`
	Errors    *string   `json:"errors"`
	Created   time.Time `json:"created"`
}
