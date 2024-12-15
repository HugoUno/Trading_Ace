package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trading-ace/internal/config"
	httpHandler "trading-ace/internal/handler/http"
	"trading-ace/internal/repository/postgres"
	"trading-ace/internal/service"
	"trading-ace/pkg/database"
	"trading-ace/pkg/ethereum"
)

func main() {
	// 載入設定
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("載入設定失敗: %v", err)
	}

	// 初始化資料庫連線
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("連接資料庫失敗: %v", err)
	}
	defer db.Close()

	// 初始化以太坊客戶端
	ethClient, err := ethereum.NewClient(cfg.Ethereum)
	if err != nil {
		log.Fatalf("建立以太坊客戶端失敗: %v", err)
	}

	// 初始化儲存層
	campaignRepo := postgres.NewCampaignRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	userTaskRepo := postgres.NewUserTaskRepository(db)

	// 初始化服務層
	campaignService := service.NewCampaignService(
		campaignRepo,
		taskRepo,
		userTaskRepo,
		ethClient,
	)

	// 初始化 HTTP 處理器
	handler := httpHandler.NewHandler(campaignService)

	// 建立 HTTP 伺服器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.Router(),
	}

	// 在背景啟動伺服器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("啟動伺服器失敗: %v", err)
		}
	}()

	// 優雅關閉
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在關閉伺服器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("強制關閉伺服器: %v", err)
	}

	log.Println("伺服器已關閉")
}
