package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"service-wallet/internal/handlers"
	"service-wallet/internal/services"
	"service-wallet/internal/store"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	walletStore, err := store.NewWalletStore(ctx)
	if err != nil {
		log.Fatal("failed init wallet store. err : ", err)
	}
	fmt.Println("🔥 Init Wallet Store...")
	defer walletStore.Close()

	walletService := services.NewWalletService(ctx, walletStore)
	if err != nil {
		log.Fatal("failed init wallet service. err : ", err)
	}
	fmt.Println("🔥 Init Wallet Service...")
	defer walletService.Close()

	walletHandler := handlers.NewUserHandler(ctx, *walletService)
	if err != nil {
		log.Fatal("failed init wallet handler. err : ", err)
	}
	fmt.Println("🔥 Init Wallet Handler...")
	defer walletHandler.Close()

	rout, err := routing(walletHandler)
	if err != nil {
		log.Fatal("failed init wallet router. err : ", err)
	}

	portHttp := os.Getenv("PORT_HTTP")
	localHost := fmt.Sprintf("0.0.0.0:%s", portHttp)
	fmt.Printf("🌐 %s\n", localHost)
	err = rout.Run(localHost)
	if err != nil {
		log.Fatalf("Could not listen and serve %s. err = %v", localHost, err)
	}
}
