package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sushanthach12/sellit-backend/internal/config"
)

func main() {
	cfg := config.MustLoad()

	args := os.Args

	// at-least 2 args required
	if len(args) < 2 {
		log.Fatal("Missing required arguments: <up | down>")
	}

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseUrl,
	)
	if err != nil {
		log.Fatalf("migrate.New: %v", err)
	}

	switch args[1] {
	case "up":
		handleUp(m)
	case "down":
		handleDown(m)
	default:
		log.Fatal("Unknown argument")
	}
}

func handleUp(m *migrate.Migrate) {
	if err := m.Up(); err != nil {
		log.Fatal("migrate.Up: ", err)
	}
}

func handleDown(m *migrate.Migrate) {
	if err := m.Down(); err != nil {
		log.Fatal("migrate.Down: ", err)
	}
}
