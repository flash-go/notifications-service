package main

import (
	// SDK
	//
	// A high-level software development toolkit based on the Flash Framework
	// for building highly efficient and fault-tolerant applications.

	"github.com/flash-go/sdk/config"
	"github.com/flash-go/sdk/infra"
	"github.com/flash-go/sdk/state"

	// Other
	"github.com/flash-go/notifications-service/internal/migrations"
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

	// Set KV from env map
	cfg.SetEnvMap(envMap)

	// Create postgres client and run migrations
	infra.NewPostgresClient(
		&infra.PostgresClientConfig{
			Cfg:        cfg,
			Telemetry:  nil,
			Migrations: migrations.Get(),
		},
	)
}
