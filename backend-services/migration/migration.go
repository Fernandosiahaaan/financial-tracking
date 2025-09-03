package main

import (
	"fmt"
	"log"
	"migration/utils"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	err = utils.CheckEnvKey([]string{
		"VERSIONING_DB",
		"POSTGRES_URI",
	})
	if err != nil {
		fmt.Println("========= INIT FAILED =========")
		log.Fatal(err)
	}
	fmt.Println("========= INIT SUCCESS =========")
}

func main() {
	Init()
	fmt.Println("========= MIGRATE PROCESS START =========")
	dbURL := os.Getenv("POSTGRES_URI")
	migrationing, err := migrate.New("file://files", dbURL)
	if err != nil {
		log.Fatalf("failed create migrate instance : %v", err)
	}

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Please provide a migration command: up | down | version <n>")
	}

	cmd := args[1]
	switch cmd {
	case "up":
		fmt.Println("========= MIGRATE UP START =========")
		err := migrationing.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("========= MIGRATE UP SUCCESS =========")

	case "down":
		fmt.Println("========= MIGRATE DOWN START =========")
		err := migrationing.Down()
		if err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		fmt.Println("========= MIGRATE DOWN SUCCESS =========")

	case "version":
		if len(args) < 3 {
			log.Fatal("Please provide a version number")
		}

		vDb, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}

		fmt.Printf("========= MIGRATE TO VERSION %d =========\n", vDb)
		err = migrationing.Migrate(uint(vDb))
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration to version failed: %v", err)
		}

		fmt.Println("========= MIGRATE TO VERSION SUCCESS =========")

	default:
		log.Fatalf("Unknown command: %s. Use up, down, or version <n>", cmd)
	}

	fmt.Println("========= MIGRATE PROCESS SUCCESSFULLY =========")
}
