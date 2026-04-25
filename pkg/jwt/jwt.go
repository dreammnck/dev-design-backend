package jwt

import (
	"backend/internal/auth"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const defaultSecret = "changeme-use-env-JWT_SECRET"

type customClaims struct {
	UserID   string        `json:"userId"`
	Username string        `json:"username"`
	Role     auth.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// jwtSecret returns the secret key from env or falls back to default
func jwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte(defaultSecret)
	}
	return []byte(secret)
}

// GenerateToken creates a signed JWT for the given user
func GenerateToken(claims auth.JWTClaims) (string, error) {
	expHours := 24 * time.Hour
	c := customClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expHours)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(jwtSecret())
}

// ValidateToken parses and validates a signed JWT, returning the embedded claims
func ValidateToken(tokenStr string) (*auth.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	c, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &auth.JWTClaims{
		UserID:   c.UserID,
		Username: c.Username,
		Role:     c.Role,
	}, nil
}
