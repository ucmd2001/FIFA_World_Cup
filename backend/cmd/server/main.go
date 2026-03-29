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
	"golang_world_cup/internal/routes"
	"golang_world_cup/internal/user"
)

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

	// 5. 啟動 Cron 排程
	scheduler.InitScheduler(db)

	// 6. 建立 Gin Engine
	r := gin.Default()

	// CORS Setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 7. 掛載路由 (依賴注入各 Handler)
	routes.SetupRouter(r, authSvc, userHdlr, matchHdlr, betHdlr)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running at :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
