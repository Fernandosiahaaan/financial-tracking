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

	// Action
	versioningDB := os.Getenv("VERSIONING_DB")
	vDb, _ := strconv.Atoi(versioningDB)
	if err != nil {
		log.Fatalf("failed convert 'VERSIONING_DB' to UINT : %s. Err : %v", versioningDB, err)
	}

	err = migrationing.Migrate(uint(vDb))
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	fmt.Println("========= MIGRATE PROCESS SUCCESSFULLY =========")
}
