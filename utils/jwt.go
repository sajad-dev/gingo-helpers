package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sajad-dev/gingo-helpers/internal/config"
)

type JWTCLeaims struct {
	Parameters map[string]any
	jwt.StandardClaims
}

func CreateJWT(field map[string]any) (string, error) {
	claims := &JWTCLeaims{
		Parameters: field,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.ConfigStUtils.JWT))
}

func ValidJWT(jwtToken string) (*JWTCLeaims, bool, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &JWTCLeaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.ConfigStUtils.JWT), nil
	})

	if err != nil {
		return nil, false, err
	}
	if claims, ok := token.Claims.(*JWTCLeaims); ok && token.Valid {
		return claims, ok, nil
	} else {
		return nil, false, nil
	}
}
