package helper

import (
	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}
