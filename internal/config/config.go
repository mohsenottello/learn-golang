// Package config reads all configuration from the environment, once, at startup.
package config

import (
	"os"

	"github.com/mohsenottello/golang/internal/mail"
)

type Config struct {
	Addr        string
	DatabaseURL string
	BaseURL     string
	SMTP        mail.Config
}

func Load() Config {
	return Config{
		Addr:        env("ADDR", ":8080"),
		DatabaseURL: mustEnv("DATABASE_URL"),
		BaseURL:     env("BASE_URL", "http://localhost:8080"),
		SMTP: mail.Config{
			Host: env("SMTP_HOST", "mailpit"),
			Port: env("SMTP_PORT", "1025"),
			User: os.Getenv("SMTP_USER"),
			Pass: os.Getenv("SMTP_PASS"),
			From: env("SMTP_FROM", "noreply@example.com"),
		},
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("missing required env var " + key)
	}
	return v
}
