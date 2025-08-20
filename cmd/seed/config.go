package main

import (
	internalConfig "github.com/flash-go/notifications-service/internal/config"
	"github.com/flash-go/sdk/infra"
	"github.com/flash-go/sdk/telemetry"
)

var envMap = map[string]string{
	"OTEL_COLLECTOR_GRPC":       telemetry.OtelCollectorGrpcOptKey,
	"OTEL_COLLECTOR_CA_CRT":     telemetry.OtelCollectorCaCrtOptKey,
	"OTEL_COLLECTOR_CLIENT_CRT": telemetry.OtelCollectorClientCrtOptKey,
	"OTEL_COLLECTOR_CLIENT_KEY": telemetry.OtelCollectorClientKeyOptKey,
	"POSTGRES_HOST":             infra.PostgresHostOptKey,
	"POSTGRES_PORT":             infra.PostgresPortOptKey,
	"POSTGRES_USER":             infra.PostgresUserOptKey,
	"POSTGRES_PASSWORD":         infra.PostgresPasswordOptKey,
	"POSTGRES_DB":               infra.PostgresDbOptKey,
	"USERS_SERVICE_NAME":        internalConfig.UsersServiceNameOptKey,
	"USERS_ADMIN_ROLE":          internalConfig.UsersAdminRoleOptKey,
	"EMAIL_SMTP_BZ_API_KEY":     internalConfig.ProvidersEmailSmtpBzApiKeyOptKey,
}
