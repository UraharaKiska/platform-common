package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret string
	AccessTTl time.Duration
	Issuer string
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Login  string `json:"login"`

	jwt.RegisteredClaims
}

func GenerateAccessToken(cfg JWTConfig, userId int64, login string) (string, error) {
	now := time.Now()
	
	claims := Claims{
		UserID: userId,
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: cfg.Issuer,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.AccessTTl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString([]byte(cfg.Secret))
}

func ParseAccessToken(cfg JWTConfig, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.Secret), nil
		},
		jwt.WithIssuer(cfg.Issuer),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

