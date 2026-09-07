package middleware

import (
	"donor-service/app/internal/config"
	"donor-service/app/internal/helper"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func ParseJwtClaims(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		token, err := echo.ContextGet[*jwt.Token](c, "user")
		if err != nil {
			return helper.Unauthorized(c, "Invalid or missing token")
		}

		claims, ok := token.Claims.(*helper.JwtCustomClaims)
		if !ok {
			return helper.Unauthorized(c, "Failed to parse claims")
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		return next(c)
	}
}

func JwtConfig() echojwt.Config {
	secret, _ := config.JWTSecret()

	return echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(helper.JwtCustomClaims)
		},
		ErrorHandler: func(c *echo.Context, err error) error {
			return helper.Unauthorized(c, "Missing or invalid token")
		},
		SigningKey: secret,
	}
}
