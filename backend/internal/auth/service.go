package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"golang_world_cup/internal/user"
)

// Service 封裝 Auth 業務邏輯
type Service struct {
	userRepo user.Repository
}

func NewService(userRepo user.Repository) *Service {
	return &Service{userRepo: userRepo}
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("super-secret-default-key-for-dev-only")
	}
	return []byte(secret)
}

// RegisterInput 定義註冊所需輸入
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput 定義登入所需輸入
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 處理使用者註冊
func (s *Service) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, _ := s.userRepo.CountByUsername(input.Username)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	total, _ := s.userRepo.CountAll()
	role := "user"
	if total == 0 {
		role = "admin"
	}

	u := &user.User{
		Username:     input.Username,
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		Points:       0,
	}
	if err := s.userRepo.Create(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

// Login 處理使用者登入
func (s *Service) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := s.userRepo.FindByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	now := time.Now()
	_ = s.userRepo.Update(u, map[string]interface{}{"last_login_at": &now})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  u.ID,
		"username": u.Username,
		"role":     u.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       u.ID,
			"username": u.Username,
			"name":     u.Name,
			"role":     u.Role,
			"points":   u.Points,
		},
	})
}

// AuthMiddleware 驗證 JWT Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("userID", claims["user_id"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

// AdminMiddleware 驗證管理員身分
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin portal access denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 確保 gorm 被使用（避免 unused import）
var _ = gorm.ErrRecordNotFound
