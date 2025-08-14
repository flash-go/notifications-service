package main

// @title		notifications-service
// @version		1.0
// @BasePath	/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"os"

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
	"github.com/flash-go/sdk/infra"
	"github.com/flash-go/sdk/logger"
	"github.com/flash-go/sdk/state"
	"github.com/flash-go/sdk/telemetry"

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
			Address:            os.Getenv("CONSUL_ADDR"),
			CAPem:              config.GetEnvBase64("CONSUL_CA_CRT"),
			CertPEM:            config.GetEnvBase64("CONSUL_CLIENT_CRT"),
			KeyPEM:             config.GetEnvBase64("CONSUL_CLIENT_KEY"),
			InsecureSkipVerify: config.GetEnvBool("CONSUL_INSECURE_SKIP_VERIFY"),
			Token:              os.Getenv("CONSUL_TOKEN"),
		},
	)

	// Create config
	cfg := config.New(
		stateService,
		os.Getenv("SERVICE_NAME"),
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

	// Add routes
	httpServer.
		// Emails

		// Send email
		AddRoute(
			http.MethodPost,
			"/emails/send",
			emailsHandler.Send,
		)

	// Register service
	if err := httpServer.RegisterService(
		os.Getenv("SERVICE_NAME"),
		os.Getenv("SERVICE_HOST"),
		config.GetEnvInt("SERVICE_PORT"),
	); err != nil {
		loggerService.Log().Err(err).Send()
	}

	// Listen http server
	if err := <-httpServer.Listen(
		os.Getenv("SERVER_HOST"),
		config.GetEnvInt("SERVER_PORT"),
	); err != nil {
		loggerService.Log().Err(err).Send()
	}
}
