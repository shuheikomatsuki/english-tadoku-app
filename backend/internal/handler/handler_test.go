package handler

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var testDB *sqlx.DB // 他のファイルで使うためにグローバル変数にしている

func TestMain(m *testing.M) {
	dsn := "host=localhost port=5432 user=postgres password=password dbname=tadoku_db sslmode=disable"
	var err error
	testDB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}
	if err := testDB.Ping(); err != nil {
		log.Fatalf("failed to ping to test database: %v", err)
	}

	exitCode := m.Run()

	testDB.Close()
	os.Exit(exitCode)
}

func cleanupUserTable(t *testing.T) {
	_, err := testDB.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("failed to cleanup users table: %v", err)
	}
}