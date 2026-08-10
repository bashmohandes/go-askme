package framework

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultPort          = 8080
	defaultSessionLife   = 20 * time.Minute
	defaultSessionCookie = ".session_id"
	defaultPostgresHost  = "localhost"
	defaultPostgresPort  = 5432
)

// Config contains the process configuration needed by the application.
type Config struct {
	Debug              bool
	Port               int
	SessionMaxLifeTime time.Duration
	SessionCookie      string
	PostgresUser       string
	PostgresPassword   string
	PostgresDB         string
	PostgresHost       string
	PostgresPort       int
	OktaClient         string
	OktaSecret         string
	OktaIssuer         string
	OktaSocialIdps     []OktaSocialIdp
}

// OktaSocialIdp is retained until the Okta authentication flow is removed.
type OktaSocialIdp struct {
	ID   string
	Name string
}

// LoadConfig reads and validates application configuration from the environment.
func LoadConfig() (*Config, error) {
	return loadConfig(os.LookupEnv)
}

type envLookup func(string) (string, bool)

func loadConfig(lookup envLookup) (*Config, error) {
	port, err := intSetting(lookup, "PORT", defaultPort)
	if err != nil {
		return nil, err
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("PORT must be between 1 and 65535")
	}

	postgresPort, err := intSetting(lookup, "POSTGRES_PORT", defaultPostgresPort)
	if err != nil {
		return nil, err
	}
	if postgresPort < 1 || postgresPort > 65535 {
		return nil, fmt.Errorf("POSTGRES_PORT must be between 1 and 65535")
	}

	debug, err := boolSetting(lookup, "DEBUG_MODE", false)
	if err != nil {
		return nil, err
	}

	sessionLife, err := durationSetting(lookup, "SESSION_MAX_LIFE_TIME", defaultSessionLife)
	if err != nil {
		return nil, err
	}
	if sessionLife <= 0 {
		return nil, fmt.Errorf("SESSION_MAX_LIFE_TIME must be positive")
	}

	config := &Config{
		Debug:              debug,
		Port:               port,
		SessionMaxLifeTime: sessionLife,
		SessionCookie:      stringSetting(lookup, "SESSION_COOKIE", defaultSessionCookie),
		PostgresUser:       stringSetting(lookup, "POSTGRES_USER", ""),
		PostgresPassword:   stringSetting(lookup, "POSTGRES_PASSWORD", ""),
		PostgresDB:         stringSetting(lookup, "POSTGRES_DB", ""),
		PostgresHost:       stringSetting(lookup, "POSTGRES_HOST", defaultPostgresHost),
		PostgresPort:       postgresPort,
		OktaClient:         stringSetting(lookup, "OKTA_CLIENT_ID", ""),
		OktaSecret:         stringSetting(lookup, "OKTA_CLIENT_SECRET", ""),
		OktaIssuer:         stringSetting(lookup, "OKTA_ISSUER", ""),
	}

	var missing []string
	for _, setting := range []struct {
		name  string
		value string
	}{
		{name: "POSTGRES_USER", value: config.PostgresUser},
		{name: "POSTGRES_PASSWORD", value: config.PostgresPassword},
		{name: "POSTGRES_DB", value: config.PostgresDB},
		{name: "OKTA_CLIENT_ID", value: config.OktaClient},
		{name: "OKTA_CLIENT_SECRET", value: config.OktaSecret},
		{name: "OKTA_ISSUER", value: config.OktaIssuer},
	} {
		if setting.value == "" {
			missing = append(missing, setting.name)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required configuration: %s", strings.Join(missing, ", "))
	}

	for _, idp := range []struct {
		key  string
		name string
	}{
		{key: "OKTA_SOCIAL_LINKEDIN_IDP", name: "LINKEDIN"},
		{key: "OKTA_SOCIAL_FACEBOOK_IDP", name: "FACEBOOK"},
	} {
		if value := stringSetting(lookup, idp.key, ""); value != "" {
			config.OktaSocialIdps = append(config.OktaSocialIdps, OktaSocialIdp{ID: value, Name: idp.name})
		}
	}

	return config, nil
}

func stringSetting(lookup envLookup, name, fallback string) string {
	value, ok := lookup(name)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func intSetting(lookup envLookup, name string, fallback int) (int, error) {
	value := stringSetting(lookup, name, "")
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", name, err)
	}
	return parsed, nil
}

func boolSetting(lookup envLookup, name string, fallback bool) (bool, error) {
	value := stringSetting(lookup, name, "")
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", name, err)
	}
	return parsed, nil
}

func durationSetting(lookup envLookup, name string, fallback time.Duration) (time.Duration, error) {
	value := stringSetting(lookup, name, "")
	if value == "" {
		return fallback, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a duration: %w", name, err)
	}
	return parsed, nil
}
