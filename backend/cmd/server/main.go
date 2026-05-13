package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"golang_world_cup/internal/auth"
	"golang_world_cup/internal/bet"
	"golang_world_cup/internal/infrastructure/database"
	"golang_world_cup/internal/infrastructure/scheduler"
	"golang_world_cup/internal/match"
	"golang_world_cup/internal/router" // 更新 import
	"golang_world_cup/internal/user"
)

// @title           FIFA World Cup Betting System API
// @version         1.0
// @description     這是一個基於 DDD 架構開發的世足賽下注系統後端 API
// @host            localhost:8080
// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 載入 .env 變數 (如有)
	_ = godotenv.Load()

	// 1. 初始化基礎設施
	db := database.InitDB("world_cup.db")

	// 2. 建立各 Domain 的 Repository
	userRepo := user.NewRepository(db)
	matchRepo := match.NewRepository(db)
	betRepo := bet.NewRepository(db)

	// 3. 建立各 Domain 的 Service
	userSvc := user.NewService(userRepo)
	matchSvc := match.NewService(matchRepo)
	betSvc := bet.NewService(betRepo)
	authSvc := auth.NewService(userRepo)

	// 4. 建立各 Domain 的 Handler
	userHdlr := user.NewHandler(userSvc)
	matchHdlr := match.NewHandler(matchSvc)
	betHdlr := bet.NewHandler(betSvc)
	authHdlr := auth.NewHandler(authSvc) // 建立 auth.Handler

	// 5. 啟動 Cron 排程
	scheduler.InitScheduler(db)

	// 6. 建立 Gin Engine
	r := gin.Default()

	// CORS Setup
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 7. 掛載路由 (依賴注入各 Handler)
	router.SetupRouter(r, authHdlr, userHdlr, matchHdlr, betHdlr) // 使用 router.SetupRouter

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running at :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
