package main

import (
	"testing"
	"time"

	"github.com/bashmohandes/go-askme/web/framework"
)

func TestBuildApplicationWiresDependencies(t *testing.T) {
	config := &framework.Config{
		Port:               8080,
		SessionMaxLifeTime: time.Minute,
		SessionCookie:      "test_session",
	}

	app, err := buildApplication(config)
	if err != nil {
		t.Fatalf("buildApplication() error = %v", err)
	}
	if app == nil {
		t.Fatal("buildApplication() = nil, want application")
	}
}
