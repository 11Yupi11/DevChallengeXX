package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var conn *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	var err error
	conn, err = sql.Open("sqlite3", "./test_db.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createTableQuery := Migrations
	_, err = conn.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func shutdown() {
	cleanup()
	conn.Close()
}

func cleanup() {
	_, _ = conn.Exec("DELETE FROM string_array")
	_, _ = conn.Exec("DELETE FROM dev_challenge")

	_, _ = conn.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'string_array'")
	_, _ = conn.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'dev_challenge'")

	_, _ = conn.Exec("INSERT INTO dev_challenge (sheet_id, cell_id, cell_value, cell_result) VALUES ('sheet0','cell0','0',0) ON CONFLICT DO NOTHING")
}
