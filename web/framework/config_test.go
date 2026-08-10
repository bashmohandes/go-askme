package framework

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestLoadConfigUsesDefaults(t *testing.T) {
	values := requiredConfig()

	config, err := loadConfig(mapLookup(values))
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}

	if config.Port != defaultPort {
		t.Errorf("Port = %d, want %d", config.Port, defaultPort)
	}
	if config.PostgresPort != defaultPostgresPort {
		t.Errorf("PostgresPort = %d, want %d", config.PostgresPort, defaultPostgresPort)
	}
	if config.Debug {
		t.Error("Debug = true, want false")
	}
	if config.PublicFolder != defaultPublicFolder {
		t.Errorf("PublicFolder = %q, want %q", config.PublicFolder, defaultPublicFolder)
	}
	if config.SessionMaxLifeTime != defaultSessionLife {
		t.Errorf("SessionMaxLifeTime = %s, want %s", config.SessionMaxLifeTime, defaultSessionLife)
	}
	if config.SessionCookie != defaultSessionCookie {
		t.Errorf("SessionCookie = %q, want %q", config.SessionCookie, defaultSessionCookie)
	}
}

func TestLoadConfigReadsOverridesAndSocialProviders(t *testing.T) {
	values := requiredConfig()
	values["PORT"] = "9090"
	values["POSTGRES_PORT"] = "5433"
	values["DEBUG_MODE"] = "true"
	values["PUBLIC_FOLDER"] = "custom/public"
	values["SESSION_MAX_LIFE_TIME"] = "2h"
	values["SESSION_COOKIE"] = "askme_session"
	values["OKTA_SOCIAL_LINKEDIN_IDP"] = "linkedin-id"
	values["OKTA_SOCIAL_FACEBOOK_IDP"] = "facebook-id"

	config, err := loadConfig(mapLookup(values))
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}

	if config.Port != 9090 || config.PostgresPort != 5433 || !config.Debug {
		t.Errorf("numeric/bool overrides not applied: %+v", config)
	}
	if config.PublicFolder != "custom/public" || config.SessionCookie != "askme_session" {
		t.Errorf("string overrides not applied: %+v", config)
	}
	if config.SessionMaxLifeTime != 2*time.Hour {
		t.Errorf("SessionMaxLifeTime = %s, want 2h", config.SessionMaxLifeTime)
	}
	wantIDPs := []OktaSocialIdp{{ID: "linkedin-id", Name: "LINKEDIN"}, {ID: "facebook-id", Name: "FACEBOOK"}}
	if !reflect.DeepEqual(config.OktaSocialIdps, wantIDPs) {
		t.Errorf("OktaSocialIdps = %#v, want %#v", config.OktaSocialIdps, wantIDPs)
	}
}

func TestLoadConfigRejectsMissingRequiredSettings(t *testing.T) {
	_, err := loadConfig(mapLookup(nil))
	if err == nil {
		t.Fatal("loadConfig() error = nil, want missing configuration error")
	}

	for _, name := range []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB", "OKTA_CLIENT_ID", "OKTA_CLIENT_SECRET", "OKTA_ISSUER"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error %q does not mention %s", err, name)
		}
	}
}

func TestLoadConfigRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{name: "non-numeric port", key: "PORT", value: "http"},
		{name: "port out of range", key: "PORT", value: "65536"},
		{name: "non-numeric postgres port", key: "POSTGRES_PORT", value: "postgres"},
		{name: "invalid debug flag", key: "DEBUG_MODE", value: "sometimes"},
		{name: "invalid session duration", key: "SESSION_MAX_LIFE_TIME", value: "tomorrow"},
		{name: "non-positive session duration", key: "SESSION_MAX_LIFE_TIME", value: "0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := requiredConfig()
			values[tt.key] = tt.value

			if _, err := loadConfig(mapLookup(values)); err == nil {
				t.Fatalf("loadConfig() error = nil for %s=%q", tt.key, tt.value)
			}
		})
	}
}

func requiredConfig() map[string]string {
	return map[string]string{
		"POSTGRES_USER":      "askme",
		"POSTGRES_PASSWORD":  "secret",
		"POSTGRES_DB":        "askme",
		"OKTA_CLIENT_ID":     "client",
		"OKTA_CLIENT_SECRET": "client-secret",
		"OKTA_ISSUER":        "https://example.test/oauth2/default",
	}
}

func mapLookup(values map[string]string) envLookup {
	return func(name string) (string, bool) {
		value, ok := values[name]
		return value, ok
	}
}
