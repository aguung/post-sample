package auth

import (
	"time"

	"post/internal/entity"
	"post/internal/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(user *entity.User) (string, error)
	GenerateRefreshToken(user *entity.User) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey     string
	expiry        int
	refreshExpiry int
}

func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{
		secretKey:     cfg.JWT.Secret,
		expiry:        cfg.JWT.Expiry,
		refreshExpiry: cfg.JWT.RefreshExpiry,
	}
}

func (j *jwtService) GenerateToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"type":    "access",
		"exp":     time.Now().Add(time.Duration(j.expiry) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) GenerateRefreshToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Duration(j.refreshExpiry) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secretKey), nil
	})
}
