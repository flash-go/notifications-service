package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/flash-go/flash/http/client"
	"github.com/flash-go/notifications-service/internal/adapter/repository/emails/model"
	emailsRepositoryAdapterPort "github.com/flash-go/notifications-service/internal/port/adapter/repository/emails"
	"gorm.io/gorm"
)

const (
	smtpBxApiBaseUrl = "https://api.smtp.bz/v1"
)

type Config struct {
	PostgresClient *gorm.DB
	HttpClient     client.Client
	SmtpBzApiKey   string
}

func New(config *Config) emailsRepositoryAdapterPort.Interface {
	return &adapter{
		postgres:     config.PostgresClient,
		httpClient:   config.HttpClient,
		smtpBzApiKey: config.SmtpBzApiKey,
	}
}

type adapter struct {
	postgres     *gorm.DB
	httpClient   client.Client
	smtpBzApiKey string
}

// Folders

func (a *adapter) CreateFolder(ctx context.Context, data emailsRepositoryAdapterPort.CreateFolderData) (*emailsRepositoryAdapterPort.FolderResult, error) {
	now := time.Now()

	// Create model
	obj := model.EmailFolder{
		ParentId:    data.ParentId,
		Name:        data.Name,
		Description: data.Description,
		SystemFlag:  data.SystemFlag,
		Updated:     time.Unix(0, now.UnixNano()),
		Created:     time.Unix(0, now.UnixNano()),
	}

	// Save folder to database
	if err := a.postgres.WithContext(ctx).Create(&obj).Error; err != nil {
		return nil, err
	}

	// Mapping model to repository
	folder := emailsRepositoryAdapterPort.FolderResult{
		Id:          obj.Id,
		ParentId:    obj.ParentId,
		Name:        obj.Name,
		Description: obj.Description,
		SystemFlag:  obj.SystemFlag,
		Updated:     obj.Updated,
		Created:     obj.Created,
	}

	return &folder, nil
}

func (a *adapter) FilterFolders(ctx context.Context, data emailsRepositoryAdapterPort.FilterFoldersData) (*[]emailsRepositoryAdapterPort.FolderResult, error) {
	// Create model
	obj := []model.EmailFolder{}

	// Create query with context
	query := a.postgres.WithContext(ctx)

	// Filter by id
	if data.Id != nil {
		query = query.Where("id IN ?", *data.Id)
	}

	// Filter by parent_id
	if data.ParentId != nil {
		var ids []uint
		hasNull := false
		for _, v := range *data.ParentId {
			if v == nil {
				hasNull = true
			} else {
				ids = append(ids, *v)
			}
		}

		if len(ids) > 0 && hasNull {
			query = query.Where("parent_id IN ? OR parent_id IS NULL", ids)
		} else if len(ids) > 0 {
			query = query.Where("parent_id IN ?", ids)
		} else if hasNull {
			query = query.Where("parent_id IS NULL")
		}
	}

	// Filter by name
	if data.Name != nil {
		query = query.Where("name IN ?", *data.Name)
	}

	// Filter by system_flag
	if data.SystemFlag != nil {
		query = query.Where("system_flag = ?", *data.SystemFlag)
	}

	// Get folders from database
	if err := query.Find(&obj).Error; err != nil {
		return nil, err
	}

	// Mapping model to repository
	folders := make([]emailsRepositoryAdapterPort.FolderResult, len(obj))
	for i, item := range obj {
		folders[i] = emailsRepositoryAdapterPort.FolderResult{
			Id:          item.Id,
			ParentId:    item.ParentId,
			Name:        item.Name,
			Description: item.Description,
			SystemFlag:  item.SystemFlag,
			Updated:     item.Updated,
			Created:     item.Created,
		}
	}

	return &folders, nil
}

func (a *adapter) DeleteFolder(ctx context.Context, id uint) error {
	// Delete folder from database
	result := a.postgres.WithContext(ctx).Delete(&model.EmailFolder{}, "id = ?", id)

	// Check errors
	if result.Error != nil {
		return result.Error
	}

	// If folder not found
	if result.RowsAffected == 0 {
		return emailsRepositoryAdapterPort.ErrFolderNotFound
	}

	return nil
}

func (a *adapter) UpdateFolder(ctx context.Context, id uint, data map[string]any) error {
	// Update folder in database
	result := a.postgres.WithContext(ctx).Model(&model.EmailFolder{}).Where("id = ?", id).Updates(data)

	// Check errors
	if result.Error != nil {
		return result.Error
	}

	// If folder not found
	if result.RowsAffected == 0 {
		return emailsRepositoryAdapterPort.ErrFolderNotFound
	}

	return nil
}

// Emails

func (a *adapter) CreateEmail(ctx context.Context, data emailsRepositoryAdapterPort.CreateEmailData) (*emailsRepositoryAdapterPort.EmailResult, error) {
	now := time.Now()

	// Create model
	obj := model.Email{
		FolderId:    data.FolderId,
		FromEmail:   data.FromEmail,
		FromName:    data.FromName,
		Subject:     data.Subject,
		Html:        data.Html,
		Text:        data.Text,
		Description: data.Description,
		SystemFlag:  data.SystemFlag,
		Updated:     time.Unix(0, now.UnixNano()),
		Created:     time.Unix(0, now.UnixNano()),
	}

	// Save email to database
	if err := a.postgres.WithContext(ctx).Create(&obj).Error; err != nil {
		return nil, err
	}

	// Mapping model to repository
	email := emailsRepositoryAdapterPort.EmailResult{
		Id:          obj.Id,
		FolderId:    obj.FolderId,
		FromEmail:   obj.FromEmail,
		FromName:    obj.FromName,
		Subject:     obj.Subject,
		Html:        obj.Html,
		Text:        obj.Text,
		Description: obj.Description,
		SystemFlag:  obj.SystemFlag,
		Updated:     obj.Updated,
		Created:     obj.Created,
	}

	return &email, nil
}

func (a *adapter) FilterEmails(ctx context.Context, data emailsRepositoryAdapterPort.FilterEmailsData) (*[]emailsRepositoryAdapterPort.EmailResult, error) {
	// Create model
	obj := []model.Email{}

	// Create query with context
	query := a.postgres.WithContext(ctx)

	// Filter by id
	if data.Id != nil {
		query = query.Where("id IN ?", *data.Id)
	}

	// Filter by folder_id
	if data.FolderId != nil {
		var ids []uint
		hasNull := false
		for _, v := range *data.FolderId {
			if v == nil {
				hasNull = true
			} else {
				ids = append(ids, *v)
			}
		}

		if len(ids) > 0 && hasNull {
			query = query.Where("folder_id IN ? OR folder_id IS NULL", ids)
		} else if len(ids) > 0 {
			query = query.Where("folder_id IN ?", ids)
		} else if hasNull {
			query = query.Where("folder_id IS NULL")
		}
	}

	// Filter by system_flag
	if data.SystemFlag != nil {
		query = query.Where("system_flag = ?", *data.SystemFlag)
	}

	// Get emails from database
	if err := query.Find(&obj).Error; err != nil {
		return nil, err
	}

	// Mapping model to repository
	emails := make([]emailsRepositoryAdapterPort.EmailResult, len(obj))
	for i, item := range obj {
		emails[i] = emailsRepositoryAdapterPort.EmailResult{
			Id:          item.Id,
			FolderId:    item.FolderId,
			FromEmail:   item.FromEmail,
			FromName:    item.FromName,
			Subject:     item.Subject,
			Html:        item.Html,
			Text:        item.Text,
			Description: item.Description,
			SystemFlag:  item.SystemFlag,
			Updated:     item.Updated,
			Created:     item.Created,
		}
	}

	return &emails, nil
}

func (a *adapter) DeleteEmail(ctx context.Context, id uint) error {
	// Delete email from database
	result := a.postgres.WithContext(ctx).Delete(&model.Email{}, "id = ?", id)

	// Check errors
	if result.Error != nil {
		return result.Error
	}

	// If email not found
	if result.RowsAffected == 0 {
		return emailsRepositoryAdapterPort.ErrEmailNotFound
	}

	return nil
}

func (a *adapter) UpdateEmail(ctx context.Context, id uint, data map[string]any) error {
	// Update email in database
	result := a.postgres.WithContext(ctx).Model(&model.Email{}).Where("id = ?", id).Updates(data)

	// Check errors
	if result.Error != nil {
		return result.Error
	}

	// If email not found
	if result.RowsAffected == 0 {
		return emailsRepositoryAdapterPort.ErrEmailNotFound
	}

	return nil
}

func (a *adapter) Send(ctx context.Context, data emailsRepositoryAdapterPort.SendData) (*emailsRepositoryAdapterPort.EmailLogResult, error) {
	// Create buffer and multipart writer
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Set fields
	if err := writer.WriteField("from", data.FromEmail); err != nil {
		return nil, err
	}
	if err := writer.WriteField("name", data.FromName); err != nil {
		return nil, err
	}
	if err := writer.WriteField("subject", data.Subject); err != nil {
		return nil, err
	}
	if err := writer.WriteField("to", data.ToEmail); err != nil {
		return nil, err
	}
	if err := writer.WriteField("html", data.Html); err != nil {
		return nil, err
	}
	if err := writer.WriteField("text", data.Text); err != nil {
		return nil, err
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close writer: %v", err)
	}

	// Send service request
	res, err := a.httpClient.Request(
		// Context
		ctx,
		// Method
		http.MethodPost,
		// URL
		smtpBxApiBaseUrl+"/smtp/send",
		// Body opt
		client.WithRequestBodyOption(
			body.Bytes(),
		),
		// Headers opts
		client.WithRequestHeadersOption(
			client.NewRequestHeader("Authorization", a.smtpBzApiKey),
			client.NewRequestHeader("Content-Type", writer.FormDataContentType()),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("service smtp.bz unavailable: %v", err)
	}

	// Create model
	obj := model.EmailLog{
		FromEmail: data.FromEmail,
		FromName:  data.FromName,
		Subject:   data.Subject,
		ToEmail:   data.ToEmail,
		Html:      data.Html,
		Text:      data.Text,
		Created:   time.Unix(0, time.Now().UnixNano()),
	}

	switch res.StatusCode() {
	case sendEmailSuccessCode:
		var response sendSuccessResponse
		if err := json.Unmarshal(res.Body(), &response); err != nil {
			return nil, fmt.Errorf("error parsing success response body: %v", err)
		}
		obj.Status = "success"
		obj.MessageId = &response.Messageid
	case sendEmailBadRequestCode:
		var response sendErrorResponse
		if err := json.Unmarshal(res.Body(), &response); err != nil {
			return nil, fmt.Errorf("error parsing bad request response body: %v", err)
		}
		var m map[string]interface{}
		if err := json.Unmarshal(response.Errors, &m); err != nil {
			return nil, fmt.Errorf("error parsing bad request response errors: %v", err)
		}
		parts := make([]string, 0, len(m))
		for k, v := range m {
			parts = append(parts, fmt.Sprintf("%s: %v", k, v))
		}
		errors := strings.Join(parts, ", ")
		obj.Status = "error"
		obj.Errors = &errors
	case sendEmailUnautorizedCode:
		message := "Unautorized"
		obj.Status = "error"
		obj.Errors = &message
	}

	// Save email to database
	if err := a.postgres.WithContext(ctx).Create(&obj).Error; err != nil {
		return nil, err
	}

	// Map model to repository results
	results := emailsRepositoryAdapterPort.EmailLogResult(obj)

	return &results, nil
}

func (a *adapter) FilterEmailLogs(ctx context.Context, data emailsRepositoryAdapterPort.FilterEmailLogsData) (*[]emailsRepositoryAdapterPort.EmailLogResult, error) {
	// Create model
	obj := []model.EmailLog{}

	// Create query with context
	query := a.postgres.WithContext(ctx)

	// Filter by id
	if data.Id != nil {
		query = query.Where("id IN ?", *data.Id)
	}

	// Filter by from_email
	if data.FromEmail != nil {
		query = query.Where("from_email IN ?", *data.FromEmail)
	}

	// Filter by from_name
	if data.FromName != nil {
		query = query.Where("from_name IN ?", *data.FromName)
	}

	// Filter by to_email
	if data.ToEmail != nil {
		query = query.Where("to_email IN ?", *data.ToEmail)
	}

	// Filter by status
	if data.Status != nil {
		query = query.Where("status IN ?", *data.Status)
	}

	// Filter by message_id
	if data.MessageId != nil {
		query = query.Where("message_id IN ?", *data.MessageId)
	}

	// Get email logs from database
	if err := query.Find(&obj).Error; err != nil {
		return nil, err
	}

	// Mapping model to repository
	logs := make([]emailsRepositoryAdapterPort.EmailLogResult, len(obj))
	for i, item := range obj {
		logs[i] = emailsRepositoryAdapterPort.EmailLogResult(item)
	}

	return &logs, nil
}
