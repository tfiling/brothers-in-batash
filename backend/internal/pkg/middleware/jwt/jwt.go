package jwt

import (
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v4"
	jtoken "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

const (
	SigningSecret          = "secret"
	ContextKey             = "token"
	TokenExpiration        = time.Hour * 6      // 6 hours
	RefreshTokenExpiration = time.Hour * 24 * 7 // 7 days
	IDClaimField           = "ID"
	ExpiryClaimField       = "exp"
	RoleClaimField         = "role"
)

func NewAuthMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(secret),
		},
		ContextKey: ContextKey,
	})
}

func GenerateToken(username string, expiration time.Duration) (string, error) {
	claims := jtoken.MapClaims{
		IDClaimField:     username,
		ExpiryClaimField: time.Now().Add(expiration).Unix(),
		//TODO - implement user roles functionality
		RoleClaimField: "admin",
	}
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SigningSecret))
	if err != nil {
		return "", errors.Wrap(err, "could not sign JWT token")
	}
	return signedToken, nil
}
