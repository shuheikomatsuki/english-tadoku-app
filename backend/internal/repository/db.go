package repository

import (
	// "database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/shuheikomatsuki/readoku/backend/internal/ssmutil"
)

func NewDBConnection() (*sqlx.DB, error) {
	sslmode := "disable" // ローカル開発用
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if env != "" && env != "local" {
		sslmode = "require"
	}

	if err := ensureDBEnv(); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		sslmode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	log.Println("connected to PostgreSQL successfully with sqlx!")

	return db, nil
}

// ensureDBEnv loads DB_* from SSM if they are not already set.
func ensureDBEnv() error {
	type pair struct {
		envKey   string
		paramEnv string
	}

	targets := []pair{
		{envKey: "DB_HOST", paramEnv: "DB_HOST_PARAM"},
		{envKey: "DB_USER", paramEnv: "DB_USER_PARAM"},
		{envKey: "DB_PASSWORD", paramEnv: "DB_PASSWORD_PARAM"},
		{envKey: "DB_NAME", paramEnv: "DB_NAME_PARAM"},
	}

	for _, t := range targets {
		if os.Getenv(t.envKey) != "" {
			continue
		}
		paramName := os.Getenv(t.paramEnv)
		if paramName == "" {
			return fmt.Errorf("%s or %s must be set", t.envKey, t.paramEnv)
		}
		value, err := ssmutil.GetParameter(paramName)
		if err != nil {
			return fmt.Errorf("fetch %s from SSM (%s): %w", t.envKey, paramName, err)
		}
		if err := os.Setenv(t.envKey, value); err != nil {
			return fmt.Errorf("set %s from SSM: %w", t.envKey, err)
		}
	}

	return nil
}
