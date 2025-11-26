package utils

import (
	"os"
	"time"

	"github.com/nerhays/prestasi_uas/app/model"

	"github.com/golang-jwt/jwt/v5"
)

type JWTCustomClaims struct {
	UserID      string             `json:"sub"`
	FullName    string             `json:"fullName"`
	Username    string             `json:"username"`
	Role        string             `json:"role"`
	Permissions []string           `json:"perms"`
	jwt.RegisteredClaims
}

func GenerateToken(user *model.User, permissions []model.Permission) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	perms := make([]string, 0, len(permissions))
	for _, p := range permissions {
		perms = append(perms, p.Name)
	}

	claims := JWTCustomClaims{
		UserID:      user.ID,
		FullName:    user.FullName,
		Username:    user.Username,
		Role:        user.Role.Name,
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
