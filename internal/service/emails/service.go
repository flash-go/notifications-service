package service

import (
	"context"
	"encoding/json"
	"strings"
	"text/template"
	"time"

	emailsRepositoryAdapterPort "github.com/flash-go/notifications-service/internal/port/adapter/repository/emails"
	emailsServicePort "github.com/flash-go/notifications-service/internal/port/service/emails"
)

type Config struct {
	EmailsRepository emailsRepositoryAdapterPort.Interface
}

func New(config *Config) emailsServicePort.Interface {
	return &service{
		config.EmailsRepository,
	}
}

type service struct {
	emailsRepository emailsRepositoryAdapterPort.Interface
}

// Folders

func (s *service) CreateFolder(ctx context.Context, data emailsServicePort.CreateFolderData) (*emailsServicePort.FolderResult, error) {
	// Check folder exist
	var parent *[]*uint
	if data.ParentId != nil {
		parent = &[]*uint{data.ParentId}
	}
	if folders, err := s.emailsRepository.FilterFolders(
		ctx,
		emailsRepositoryAdapterPort.FilterFoldersData{
			ParentId: parent,
			Name:     &[]string{data.Name},
		},
	); err != nil {
		return nil, err
	} else if len(*folders) == 1 {
		return nil, emailsServicePort.ErrFolderExist
	}

	// Create folder
	folder, err := s.emailsRepository.CreateFolder(
		ctx,
		emailsRepositoryAdapterPort.CreateFolderData(data),
	)
	if err != nil {
		return nil, err
	}

	// Map repository to service results
	results := emailsServicePort.FolderResult(*folder)

	return &results, nil
}

func (s *service) FilterFolders(ctx context.Context, data emailsServicePort.FilterFoldersData) (*[]emailsServicePort.FolderResult, error) {
	// Filter folders
	folders, err := s.emailsRepository.FilterFolders(
		ctx,
		emailsRepositoryAdapterPort.FilterFoldersData(data),
	)
	if err != nil {
		return nil, err
	}

	// Map repository to service results
	results := make([]emailsServicePort.FolderResult, 0, len(*folders))
	for _, folder := range *folders {
		results = append(
			results,
			emailsServicePort.FolderResult(folder),
		)
	}

	return &results, nil
}

func (s *service) DeleteFolder(ctx context.Context, id uint) error {
	// Delete folder
	if err := s.emailsRepository.DeleteFolder(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateFolder(ctx context.Context, id uint, data map[string]any) error {
	// Set folder data
	data["updated"] = time.Unix(0, time.Now().UnixNano())

	// Update folder
	return s.emailsRepository.UpdateFolder(ctx, id, data)
}

// Emails

func (s *service) CreateEmail(ctx context.Context, data emailsServicePort.CreateEmailData) (*emailsServicePort.EmailResult, error) {
	// Create email
	email, err := s.emailsRepository.CreateEmail(
		ctx,
		emailsRepositoryAdapterPort.CreateEmailData(data),
	)
	if err != nil {
		return nil, err
	}

	// Map repository to service results
	results := emailsServicePort.EmailResult(*email)

	return &results, nil
}

func (s *service) FilterEmails(ctx context.Context, data emailsServicePort.FilterEmailsData) (*[]emailsServicePort.EmailResult, error) {
	// Filter emails
	emails, err := s.emailsRepository.FilterEmails(
		ctx,
		emailsRepositoryAdapterPort.FilterEmailsData(data),
	)
	if err != nil {
		return nil, err
	}

	// Map repository to service results
	results := make([]emailsServicePort.EmailResult, 0, len(*emails))
	for _, email := range *emails {
		results = append(
			results,
			emailsServicePort.EmailResult(email),
		)
	}

	return &results, nil
}

func (s *service) DeleteEmail(ctx context.Context, id uint) error {
	// Delete email
	if err := s.emailsRepository.DeleteEmail(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateEmail(ctx context.Context, id uint, data map[string]any) error {
	// Set folder data
	data["updated"] = time.Unix(0, time.Now().UnixNano())

	// Update email
	return s.emailsRepository.UpdateEmail(ctx, id, data)
}

func (s *service) SendCustom(ctx context.Context, data emailsServicePort.SendCustomData) (*emailsServicePort.EmailLogResult, error) {
	// Send email
	log, err := s.emailsRepository.Send(
		ctx,
		emailsRepositoryAdapterPort.SendData(data),
	)
	if err != nil {
		return nil, err
	}

	// Map repository to service results
	results := emailsServicePort.EmailLogResult(*log)

	return &results, nil
}

func (s *service) Send(ctx context.Context, data emailsServicePort.SendData) (*emailsServicePort.EmailLogResult, error) {
	// Get email
	emails, err := s.emailsRepository.FilterEmails(
		ctx,
		emailsRepositoryAdapterPort.FilterEmailsData{
			Id: &[]uint{data.EmailId},
		},
	)
	if err != nil {
		return nil, err
	}
	if len(*emails) == 0 {
		return nil, emailsServicePort.ErrEmailNotFound
	}

	// Render template
	subject, err := s.renderTemplate((*emails)[0].Subject, data.Vars)
	if err != nil {
		return nil, err
	}
	html, err := s.renderTemplate((*emails)[0].Html, data.Vars)
	if err != nil {
		return nil, err
	}
	text, err := s.renderTemplate((*emails)[0].Text, data.Vars)
	if err != nil {
		return nil, err
	}

	// Send email
	log, err := s.emailsRepository.Send(
		ctx,
		emailsRepositoryAdapterPort.SendData{
			FromEmail: (*emails)[0].FromEmail,
			FromName:  (*emails)[0].FromName,
			Subject:   *subject,
			ToEmail:   data.ToEmail,
			Html:      *html,
			Text:      *text,
		},
	)
	if err != nil {
		return nil, err
	}

	// Map repository to service results
	results := emailsServicePort.EmailLogResult(*log)

	return &results, nil
}

func (s *service) FilterEmailLogs(ctx context.Context, data emailsServicePort.FilterEmailLogsData) (*[]emailsServicePort.EmailLogResult, error) {
	// Filter email logs
	logs, err := s.emailsRepository.FilterEmailLogs(
		ctx,
		emailsRepositoryAdapterPort.FilterEmailLogsData(data),
	)
	if err != nil {
		return nil, err
	}

	// Map repository to service results
	results := make([]emailsServicePort.EmailLogResult, 0, len(*logs))
	for _, log := range *logs {
		results = append(
			results,
			emailsServicePort.EmailLogResult(log),
		)
	}

	return &results, nil
}

func (s *service) renderTemplate(templateContent string, vars *json.RawMessage) (*string, error) {
	// No vars
	if vars == nil {
		return &templateContent, nil
	}

	// JSON to map
	var data map[string]interface{}
	if err := json.Unmarshal(*vars, &data); err != nil {
		return nil, err
	}

	// Create template
	tmpl, err := template.New("tmpl").Parse(templateContent)
	if err != nil {
		return nil, err
	}

	// Render template
	var sb strings.Builder
	if err := tmpl.Execute(&sb, data); err != nil {
		return nil, err
	}

	result := sb.String()
	return &result, nil
}
