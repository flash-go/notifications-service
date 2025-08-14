package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/flash-go/flash/http/client"
	"github.com/flash-go/notifications-service/internal/adapter/repository/emails/model"
	"github.com/flash-go/notifications-service/internal/domain/entity"
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

func (a *adapter) Send(ctx context.Context, email *entity.Email) error {
	// Create buffer and multipart writer
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Set fields
	if err := writer.WriteField("from", email.FromEmail); err != nil {
		return err
	}
	if err := writer.WriteField("name", email.FromName); err != nil {
		return err
	}
	if err := writer.WriteField("subject", email.Subject); err != nil {
		return err
	}
	if err := writer.WriteField("to", email.ToEmail); err != nil {
		return err
	}
	if err := writer.WriteField("html", email.Html); err != nil {
		return err
	}
	if err := writer.WriteField("text", email.Text); err != nil {
		return err
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close writer: %v", err)
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
		return fmt.Errorf("service smtp.bz unavailable: %v", err)
	}

	switch res.StatusCode() {
	case sendEmailSuccessCode:
		var response sendSuccessResponse
		if err := json.Unmarshal(res.Body(), &response); err != nil {
			return fmt.Errorf("error parsing success response body: %v", err)
		}
		email.Status = "success"
		email.MessageId = &response.Messageid
	case sendEmailBadRequestCode:
		var response sendErrorResponse
		if err := json.Unmarshal(res.Body(), &response); err != nil {
			return fmt.Errorf("error parsing bad request response body: %v", err)
		}
		var m map[string]interface{}
		if err := json.Unmarshal(response.Errors, &m); err != nil {
			return fmt.Errorf("error parsing bad request response errors: %v", err)
		}
		parts := make([]string, 0, len(m))
		for k, v := range m {
			parts = append(parts, fmt.Sprintf("%s: %v", k, v))
		}
		errors := strings.Join(parts, ", ")
		email.Status = "error"
		email.Errors = &errors
	case sendEmailUnautorizedCode:
		message := "Unautorized"
		email.Status = "error"
		email.Errors = &message
	}

	// Mapping entity to model
	obj := model.Email(*email)

	// Save email to database
	if err := a.postgres.WithContext(ctx).Create(&obj).Error; err != nil {
		return err
	}

	// Mapping model to entity
	*email = entity.Email(obj)

	return nil
}
