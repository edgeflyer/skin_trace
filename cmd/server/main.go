package main

import (
	"log"

	"skintrace/internal/agent"
	"skintrace/internal/config"
	"skintrace/internal/db"
	"skintrace/internal/deepseek"
	"skintrace/internal/record"
	"skintrace/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.OpenMySQL(cfg.MySQL.DSN)
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.AutoMigrate(&record.Record{}); err != nil {
		log.Fatalf("migrate skin records: %v", err)
	}

	files, err := storage.NewLocal("uploads")
	if err != nil {
		log.Fatal(err)
	}

	repo := record.NewRepository(conn)
	service := record.NewService(repo, files)
	handler := record.NewHandler(service)

	modelClient, err := deepseek.NewClient(
		cfg.DeepSeek.BaseURL,
		cfg.DeepSeek.APIKey,
		cfg.DeepSeek.Model,
	)
	if err != nil {
		log.Fatal(err)
	}

	chatAgent := agent.NewAgent(modelClient, service)
	chatHandler := agent.NewHandler(chatAgent)

	router := gin.Default()
	router.POST("/records", handler.Create)
	router.GET("/records", handler.List)
	router.POST("/chat", chatHandler.Chat)

	router.StaticFS("/uploads", gin.Dir("uploads", true))

	if err := router.Run("127.0.0.1:8080"); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
