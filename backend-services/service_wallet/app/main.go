package main

import (
	"fmt"
	"log"
	"os"
	"service-wallet/utils"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	err = utils.CheckEnvKey([]string{
		"PORT_HTTP",
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
	fmt.Println("========= APP START =========")
	fmt.Println("test")
	rout, err := routing()
	if err != nil {
		log.Fatal(err)
	}

	portHttp := os.Getenv("PORT_HTTP")
	localHost := fmt.Sprintf("0.0.0.0:%s", portHttp)
	fmt.Printf("🌐 %s\n", localHost)
	err = rout.Run(localHost)
	if err != nil {
		log.Fatalf("Could not listen and serve %s. err = %v", localHost, err)
	}
}
