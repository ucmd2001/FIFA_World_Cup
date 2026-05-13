package router

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "golang_world_cup/docs" // 引入產生的 Swagger docs
	"golang_world_cup/internal/auth"
	"golang_world_cup/internal/bet"
	"golang_world_cup/internal/match"
	"golang_world_cup/internal/user"
)

// SetupRouter 掛載所有路由，接受各 Domain Handler 作為參數（DI）
func SetupRouter(
	r *gin.Engine,
	authHdlr *auth.Handler,
	userHdlr *user.Handler,
	matchHdlr *match.Handler,
	betHdlr *bet.Handler,
) {
	// Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")

	// 1. 公開路由 (無須驗證)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHdlr.Register)
		authGroup.POST("/login", authHdlr.Login)
	}

	api.GET("/leaderboard", userHdlr.GetLeaderboard)
	api.GET("/matches", matchHdlr.GetMatches)

	// 2. 保護路由 (需要 JWT)
	protected := api.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		userGroup := protected.Group("/user")
		{
			userGroup.GET("/me", userHdlr.GetMe)
			userGroup.POST("/claim-daily", userHdlr.ClaimDaily)
			userGroup.GET("/bets", betHdlr.GetUserBets)
		}

		betsGroup := protected.Group("/bets")
		{
			betsGroup.POST("/", betHdlr.PlaceBet)
		}
	}

	// 3. 管理員路由
	admin := api.Group("/admin")
	admin.Use(auth.AuthMiddleware(), auth.AdminMiddleware())
	{
		admin.POST("/matches", matchHdlr.CreateMatch)
		admin.POST("/matches/sync", matchHdlr.SyncMatches)
		admin.POST("/matches/:id/resolve", matchHdlr.ResolveMatch)
		admin.GET("/users", userHdlr.GetAllUsers)
		admin.PUT("/users/:id", userHdlr.UpdateUser)
		admin.DELETE("/users/:id", userHdlr.DeleteUser)
	}
}
