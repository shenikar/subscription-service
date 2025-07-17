package main

import (
	"context"
	"log"

	"github.com/shenikar/subscription-service/internal/config"
	"github.com/shenikar/subscription-service/internal/db"
	"github.com/shenikar/subscription-service/internal/handler"
	"github.com/shenikar/subscription-service/internal/repository"
	"github.com/shenikar/subscription-service/internal/router"
	"github.com/shenikar/subscription-service/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	conn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer conn.Close(context.Background())

	repo := repository.NewSubscriptionRepository(conn)
	svc := service.NewSubscriptionService(repo)
	handl := handler.NewSubscriptionHandler(svc)

	router := router.SetupRouter(handl)

	addr := ":" + cfg.ServerPort
	log.Printf("server is running at %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
