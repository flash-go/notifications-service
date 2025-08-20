package main

// @title		notifications-service
// @version		1.0
// @BasePath	/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	// Framework
	//
	// Core of the Flash Framework. Contains the fundamental components of
	// the application.

	"github.com/flash-go/flash/http"
	"github.com/flash-go/flash/http/client"
	"github.com/flash-go/flash/http/server"

	// SDK
	//
	// A high-level software development toolkit based on the Flash Framework
	// for building highly efficient and fault-tolerant applications.

	"github.com/flash-go/sdk/config"
	"github.com/flash-go/sdk/errors"
	"github.com/flash-go/sdk/http/server/middleware"
	"github.com/flash-go/sdk/infra"
	"github.com/flash-go/sdk/logger"
	"github.com/flash-go/sdk/services/users"
	"github.com/flash-go/sdk/state"
	"github.com/flash-go/sdk/telemetry"

	// Ports

	//// Handlers
	httpEmailsHandlerAdapterPort "github.com/flash-go/notifications-service/internal/port/adapter/handler/emails/http"

	// Implementations

	//// Handlers
	httpEmailsHandlerAdapterImpl "github.com/flash-go/notifications-service/internal/adapter/handler/emails/http"

	//// Repository
	emailsRepositoryAdapterImpl "github.com/flash-go/notifications-service/internal/adapter/repository/emails"

	//// Services
	emailsServiceImpl "github.com/flash-go/notifications-service/internal/service/emails"

	// Config
	internalConfig "github.com/flash-go/notifications-service/internal/config"

	// Other
	_ "github.com/flash-go/notifications-service/docs"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Create state service
	stateService := state.NewWithSecureAuth(
		&state.SecureAuthConfig{
			Address:            config.GetEnvStr("CONSUL_ADDR"),
			CAPem:              config.GetEnvBase64("CONSUL_CA_CRT"),
			CertPEM:            config.GetEnvBase64("CONSUL_CLIENT_CRT"),
			KeyPEM:             config.GetEnvBase64("CONSUL_CLIENT_KEY"),
			InsecureSkipVerify: config.GetEnvBool("CONSUL_INSECURE_SKIP_VERIFY"),
			Token:              config.GetEnvStr("CONSUL_TOKEN"),
		},
	)

	// Create config
	cfg := config.New(
		stateService,
		config.GetEnvStr("SERVICE_NAME"),
	)

	// Create logger service
	loggerService := logger.NewConsole()

	// Set log level
	loggerService.SetLevel(config.GetEnvInt("LOG_LEVEL"))

	// Create telemetry service
	telemetryService := telemetry.NewSecureGrpc(cfg)

	// Collect metrics
	telemetryService.CollectGoRuntimeMetrics(collectGoRuntimeMetricsTimeout)

	// Create postgres client without migrations
	postgresClient := infra.NewPostgresClient(
		&infra.PostgresClientConfig{
			Cfg:        cfg,
			Telemetry:  telemetryService,
			Migrations: nil,
		},
	)

	// Create http client
	httpClient := client.New()

	// Use state service
	httpClient.UseState(stateService)

	// Use telemetry service
	httpClient.UseTelemetry(telemetryService)

	// Create http server
	httpServer := server.New()

	// Use telemetry service
	httpServer.UseTelemetry(telemetryService)

	// Use logger service
	httpServer.UseLogger(loggerService)

	// Use state service
	httpServer.UseState(stateService)

	// Use Swagger
	httpServer.UseSwagger()

	// Set error response status map
	httpServer.SetErrorResponseStatusMap(
		&server.ErrorResponseStatusMap{
			errors.ErrBadRequest:   400,
			errors.ErrUnauthorized: 401,
			errors.ErrForbidden:    403,
			errors.ErrNotFound:     404,
		},
	)

	// Create repository
	emailsRepository := emailsRepositoryAdapterImpl.New(
		&emailsRepositoryAdapterImpl.Config{
			PostgresClient: postgresClient,
			HttpClient:     httpClient,
			SmtpBzApiKey:   cfg.Get(internalConfig.ProvidersEmailSmtpBzApiKeyOptKey),
		},
	)

	// Create services
	emailsService := emailsServiceImpl.New(
		&emailsServiceImpl.Config{
			EmailsRepository: emailsRepository,
		},
	)

	// Create handlers
	emailsHandler := httpEmailsHandlerAdapterImpl.New(
		&httpEmailsHandlerAdapterImpl.Config{
			EmailsService: emailsService,
		},
	)

	// Create users middleware
	usersMiddleware := users.NewMiddleware(
		&users.MiddlewareConfig{
			UsersService: cfg.Get(internalConfig.UsersServiceNameOptKey),
			HttpClient:   httpClient,
		},
	)

	// Get admin role
	adminRole := cfg.Get(internalConfig.UsersAdminRoleOptKey)

	// Add routes
	httpServer.
		// Emails

		// Create email folder (admin)
		AddRoute(
			http.MethodPost,
			"/admin/notifications/emails/folders",
			emailsHandler.AdminCreateFolder,
			middleware.ParseJsonBody[*httpEmailsHandlerAdapterPort.CreateFolderData](),
			usersMiddleware.Auth(
				users.WithAuthRolesOption(adminRole),
			),
		).
		// Filter email folders (admin)
		AddRoute(
			http.MethodPost,
			"/admin/notifications/emails/folders/filter",
			emailsHandler.AdminFilterFolders,
			middleware.ParseJsonBody[*httpEmailsHandlerAdapterPort.FilterFoldersData](),
			usersMiddleware.Auth(
				users.WithAuthRolesOption(adminRole),
			),
		).
		// Delete email folder (admin)
		AddRoute(
			http.MethodDelete,
			"/admin/notifications/emails/folders/{id}",
			emailsHandler.AdminDeleteFolder,
			usersMiddleware.Auth(
				users.WithAuthRolesOption(adminRole),
			),
		).
		// Update email folder (admin)
		AddRoute(
			http.MethodPatch,
			"/admin/notifications/emails/folders/{id}",
			emailsHandler.AdminUpdateFolder,
			middleware.ParseJsonBody[*httpEmailsHandlerAdapterPort.UpdateFolderData](),
			usersMiddleware.Auth(
				users.WithAuthRolesOption(adminRole),
			),
		).

		// Create email (admin)
		AddRoute(
			http.MethodPost,
			"/admin/notifications/emails",
			emailsHandler.AdminCreateEmail,
			middleware.ParseJsonBody[*httpEmailsHandlerAdapterPort.CreateEmailData](),
			usersMiddleware.Auth(
				users.WithAuthRolesOption(adminRole),
			),
		).
		// Filter emails (admin)
		AddRoute(
			http.MethodPost,
			"/admin/notifications/emails/filter",
			emailsHandler.AdminFilterEmails,
			middleware.ParseJsonBody[*httpEmailsHandlerAdapterPort.FilterEmailsData](),
			usersMiddleware.Auth(
				users.WithAuthRolesOption(adminRole),
			),
		).
		// Delete email (admin)
		AddRoute(
			http.MethodDelete,
			"/admin/notifications/emails/{id}",
			emailsHandler.AdminDeleteEmail,
			usersMiddleware.Auth(
				users.WithAuthRolesOption(adminRole),
			),
		).
		// Update email (admin)
		AddRoute(
			http.MethodPatch,
			"/admin/notifications/emails/{id}",
			emailsHandler.AdminUpdateEmail,
			middleware.ParseJsonBody[*httpEmailsHandlerAdapterPort.UpdateEmailData](),
			usersMiddleware.Auth(
				users.WithAuthRolesOption(adminRole),
			),
		).

		// Send custom email
		AddRoute(
			http.MethodPost,
			"/notifications/emails/send/custom",
			emailsHandler.SendCustom,
			middleware.ParseJsonBody[*httpEmailsHandlerAdapterPort.SendCustomData](),
			usersMiddleware.Auth(
				users.WithAuthRolesOption(adminRole),
			),
		).
		// Send email
		AddRoute(
			http.MethodPost,
			"/notifications/emails/send",
			emailsHandler.Send,
			middleware.ParseJsonBody[*httpEmailsHandlerAdapterPort.SendData](),
		).
		// Filter email logs (admin)
		AddRoute(
			http.MethodPost,
			"/admin/notifications/emails/logs/filter",
			emailsHandler.AdminFilterEmailLogs,
			middleware.ParseJsonBody[*httpEmailsHandlerAdapterPort.FilterEmailLogsData](),
			usersMiddleware.Auth(
				users.WithAuthRolesOption(adminRole),
			),
		)

	// Register service
	if err := httpServer.RegisterService(
		config.GetEnvStr("SERVICE_NAME"),
		config.GetEnvStr("SERVICE_HOST"),
		config.GetEnvInt("SERVICE_PORT"),
	); err != nil {
		loggerService.Log().Err(err).Send()
	}

	// Listen http server
	if err := <-httpServer.Listen(
		config.GetEnvStr("SERVER_HOST"),
		config.GetEnvInt("SERVER_PORT"),
	); err != nil {
		loggerService.Log().Err(err).Send()
	}
}
