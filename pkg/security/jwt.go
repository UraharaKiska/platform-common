package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret string `toml:"secret,omitempty"`
	AccessTTl time.Duration `toml:"access_ttl,omitempty"`
	Issuer string `toml:"issuer,omitempty"`
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Login  string `json:"login"`
	Roles []string `json:"roles"`

	jwt.RegisteredClaims
}

func GenerateAccessToken(cfg JWTConfig, userId int64, login string, roles []string) (string, error) {
	now := time.Now()
	
	claims := Claims{
		UserID: userId,
		Login: login,
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: cfg.Issuer,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.AccessTTl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

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

