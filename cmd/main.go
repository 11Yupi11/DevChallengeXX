package main

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/jinzhu/configor"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"dev-challenge/db"
	"dev-challenge/internal/config"
	"dev-challenge/internal/server"
)

var (
	cfg config.Config
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "etc/config.yml"
	}

	if err := configor.Load(&cfg, configPath); err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	// Override with environment variables, if they're set
	if port, exists := os.LookupEnv("APP_PORT"); exists {
		cfg.Port, _ = strconv.Atoi(port)
	}

	if debug, exists := os.LookupEnv("APP_DEBUG"); exists {
		cfg.Debug = debug == "true"
	}

	store := mustOpenDBConnection()

	s, err := server.NewServer(store, &cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	s.Run()
}

func mustOpenDBConnection() *sql.DB {
	database, _ := sql.Open("sqlite3", "./persistent_storage/main.db")

	createTableQuery := db.Migrations
	_, err := database.Exec(createTableQuery)
	if err != nil {
		logrus.Fatalf("Failed to create table: %v", err)
	}
	return database
}
