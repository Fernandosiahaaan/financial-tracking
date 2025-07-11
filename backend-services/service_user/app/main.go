package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"service-user/infrastructure/redis"
	"service-user/internal/handler"
	"service-user/internal/service"
	"service-user/internal/store"
	"service-user/middlewares"
	"service-user/utils"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	err = utils.CheckEnvKey([]string{
		"PORT_HTTP",
		"REDIS_HOST",
		"REDIS_PORT",
		"POSTGRES_URI",
		"SECRET_KEY",
	})
	if err != nil {
		fmt.Println("========= INIT FAILED =========")
		log.Fatal(err)
	}
	fmt.Println("========= INIT SUCCESS =========")
}

func main() {
	Init()
	log.Println("========= APP START =========")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userStore, err := store.NewUserStore(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer userStore.Close()
	fmt.Println("🔥 Init Repository...")

	userRedis, err := redis.NewReddisClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer userRedis.Close()
	fmt.Println("🔥 Init Redis...")

	userService := service.NewUserService(ctx, userRedis, userStore)
	defer userService.Close()
	fmt.Println("🔥 Init Service...")

	mw := middlewares.NewMidleware(ctx)
	defer mw.Close()
	fmt.Println("🔥 Init midleware...")

	userHandler := handler.NewUserHandler(handler.ParamHandler{
		Ctx:     ctx,
		Service: userService,
		Redis:   userRedis,
		Store:   userStore,
	})
	defer userHandler.Close()
	fmt.Println("🔥 Init Handler...")

	// Routing
	rout, err := routing(userHandler)
	if err != nil {
		log.Fatal(err)
	}

	// HTTP Start
	portHttp := os.Getenv("PORT_HTTP")
	localHost := fmt.Sprintf("0.0.0.0:%s", portHttp)
	fmt.Printf("🌐 %s\n", localHost)
	err = rout.Run(localHost)
	if err != nil {
		log.Fatalf("Could not listen and serve %s. err = %v", localHost, err)
	}
}
