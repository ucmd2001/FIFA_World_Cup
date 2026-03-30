package auth

import (
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

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

// Register 處理使用者註冊邏輯
func (s *Service) Register(input RegisterInput) error {
	count, _ := s.userRepo.CountByUsername(input.Username)
	if count > 0 {
		return errors.New("Username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Failed to hash password")
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
		return errors.New("Failed to create user")
	}

	return nil
}

// Login 處理使用者登入並回傳結果 (此處目前暫留 gin.Context 以利漸進重構，但邏輯已分離)
func (s *Service) Login(c *gin.Context, input LoginInput) error {
	u, err := s.userRepo.FindByUsername(input.Username)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return err
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
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return err
	}

	c.JSON(200, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       u.ID,
			"username": u.Username,
			"name":     u.Name,
			"role":     u.Role,
			"points":   u.Points,
		},
	})
	return nil
}
